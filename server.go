package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"

	"ip-calc-practice-api/db"
	"ip-calc-practice-api/models"
)

const (
	COOKIE_NAME string = "icp-id"
	LIMIT_TOTAL int = 100
	DATETIME_FORMAT string = "2006-01-02T15:04:05"
)

var repo *db.IpRepository

// summary => 待ち受けるサーバのルータを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("api/v1")

	v1.GET("/init/:total", initializeAction)
	v1.POST("/next/:id/:number", getNextQuestion)
	//v1.GET("/resume/:id", resumeAnswer)
	//v1.DELETE("clean", deleteExpiredData)

	repo = initDB()

	 return router
}

// summary => 最初から始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func initializeAction(c *gin.Context) {
	total_prm := c.Param("total")
	total := getTotalValue(total_prm)

	id := getID()
	tq := generateNewQuestion(id, 1)

	mid := models.MstrID{
		Id:     id,
		Total:  total,
		Expire: time.Now().AddDate(0, 0, 1).Format(DATETIME_FORMAT),
	}

	if repo == nil {
		c.JSON(http.StatusServiceUnavailable, errCannotConnectDB)
		c.Abort()
		return
	}

	if err := repo.InsertFirstData(mid, tq); err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedOperateData)
		c.Abort()
		return
	}

	c.SetCookie(COOKIE_NAME, id, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, getQuestionSet(tq))
}

// summary => 解答をDBへ格納して、次の問題 or 結果を返す処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func getNextQuestion(c *gin.Context) {
	id_prm  := c.Param("id")
	num_prm := c.Param("number")

	number, err := strconv.Atoi(num_prm)
	if err != nil || (number < 1 || number > LIMIT_TOTAL) {
		c.JSON(http.StatusBadRequest, errInvalidRequestedURL)
		c.Abort()
		return
	}

	if repo == nil {
		c.JSON(http.StatusServiceUnavailable, errCannotConnectDB)
		c.Abort()
		return
	}

	check_qnum, err := repo.CheckNow(id_prm)

	if err != nil {
		c.JSON(http.StatusBadRequest, errFailedGetData)
		c.Abort()
		return
	}

	if number != check_qnum.Now {
		c.JSON(http.StatusBadRequest, errInvalidRequestedData)
		c.Abort()
		return
	}

	tq, err := repo.GetQuestion(id_prm, number)

	if err != nil {
		c.JSON(http.StatusBadRequest, errFailedGetData)
		c.Abort()
		return
	}

	var req models.AnswerSet

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, errInvalidRequestedData)
		c.Abort()
		return
	}

	nwaddr := concatIPAddress(req.NwAddr1st, req.NwAddr2nd,
		                      req.NwAddr3rd, req.NwAddr4th)
	bcaddr := concatIPAddress(req.BcAddr1st, req.BcAddr2nd,
		                      req.BcAddr3rd, req.BcAddr4th)
	elapsed, _ := strconv.Atoi(req.Elapsed)

	tq.AnsNwAddr = nwaddr
	tq.AnsBcAddr = bcaddr
	tq.Elapsed   = elapsed
	tq.Updated   = getNowString()

	if err := repo.UpdateQuestion(tq); err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedOperateData)
		c.Abort()
		return
	}

	if number < check_qnum.Total {
		newq := generateNewQuestion(id_prm, check_qnum.Now + 1)

		if err := repo.InsertQuestion(newq); err != nil {
			c.JSON(http.StatusServiceUnavailable, errFailedOperateData)
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, getQuestionSet(newq))

	} else {
		results, err := repo.GetResults(id_prm)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, errFailedGetData)
			c.Abort()
			return
		}

		rc := getResultCollection(results)
		rc.IsEnd = true

		c.JSON(http.StatusOK, rc)
	}
}

/*
// summary => 続きから始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func resumeAnswer(c *gin.Context) {
	_, err := c.Cookie(COOKIE_NAME)
	if err == nil {
		c.SetCookie(COOKIE_NAME, "", -1, "/", "", false, true)
	}

	m, err := repo.GetExpired(id)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedGetData)
		c.Abort()
		return
	}
}

func deleteExpiredData(c *gin.Context) {
}
*/

// summary => DBとの接続についての初期処理
// return::*db.IpRepository => 構造体
/////////////////////////////////////////
func initDB() *db.IpRepository {
	driver, dsn, err := db.GetDataSourceName()
	if err != nil {
		fmt.Println("E001 :", err)
		return nil
	}

	var dbmap *gorp.DbMap

	switch driver {
	case "mysql":
		op, _ := sql.Open(driver, dsn)
		dial  := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}

		dbmap = &gorp.DbMap{Db: op, Dialect: dial, ExpandSliceArgs: true}
		models.MapStructsToTables(dbmap)
	}

	return db.NewIpRepository(dbmap)
}

// summary => 固有IDを取得します
// return::string => 一意な固有ID
/////////////////////////////////////////
func getID() string {
	obj, err := uuid.NewRandom()

	if err == nil {
		return obj.String()
	} else {
		return uuid.Nil.String()
	}
}

// summary => 変換元文字列から出題する問題数を決定します
// param::totalStr => Paramからの流入値
// return::int => 出題数
/////////////////////////////////////////
func getTotalValue(totalStr string) int {
	total, err := strconv.Atoi(totalStr)

	// 変換に失敗 OR :totalのパラメータが存在しない場合、10問
	if err != nil {
		return 10
	}

	if total < 1 {
		// 0以下の数値が入っている場合、10問
		total = 10

	} else if total > LIMIT_TOTAL {
		// 上限値より大きい値が入っている場合、強制的に上限値
		total = LIMIT_TOTAL
	}

	return total
}

// summary => 新しい問題を生成し、その構造体を返します
/////////////////////////////////////////
func generateNewQuestion(id string, num int) models.TranQuestion {
	ipint, bits := getQuestion()
	ip  := uint2ip(ipint)
	now := getNowString()

	return models.TranQuestion{
		Id       : id,
		Number   : num,
		Source   : ip.String(),
		CIDRbits : bits,
		IsCIDR   : int(ipint % 2),
		CorNwAddr: getNetworkAddress(ip, bits).String(),
		AnsNwAddr: "0.0.0.0",
		CorBcAddr: getBroadcastAddress(ip, bits).String(),
		AnsBcAddr: "0.0.0.0",
		Elapsed  : 0,
		Created  : now,
		Updated  : now,
	}
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getQuestionSet(tq models.TranQuestion) models.QuestionSet {
	qs := models.QuestionSet{
		Id        : tq.Id,
		Number    : tq.Number,
		Source    : tq.Source,
		CIDRbits  : tq.CIDRbits,
		SubnetMask: "",
	}

	if tq.IsCIDR == 1 {
		qs.CIDRbits   = -1
		qs.SubnetMask = getSubnetMask(tq.CIDRbits)
	}

	return qs
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getResultCollection(tqList []models.TranQuestion) models.ResultCollection {
	r   := make([]models.ResultSet, len(tqList))
	res := models.ResultCollection{
		Id:     "",
		IsEnd:  false,
		Result: r,
	}

	for i, tq := range tqList {
		rs := models.ResultSet{
			Number     : tq.Number,
			Source     : tq.Source,
			CIDRbits   : tq.CIDRbits,
			SubnetMask : "",
			CorNwAddr  : tq.CorNwAddr,
			AnsNwAddr  : tq.AnsNwAddr,
			CorBcAddr  : tq.CorBcAddr,
			AnsBcAddr  : tq.AnsBcAddr,
			AnswerdTime: 0,
		}

		if tq.IsCIDR == 1 {
			rs.CIDRbits   = -1
			rs.SubnetMask = getSubnetMask(tq.CIDRbits)
		}

		if i == 0 {
			res.Id = tq.Id
			rs.AnswerdTime = tq.Elapsed

		} else {
			if tq.Elapsed != 0 {
				rs.AnswerdTime = tq.Elapsed - tqList[i-1].Elapsed
			}
		}

		res.Result[i] = rs
	}

	return res
}

// summary => SQLに登録されている文字列型の時間をtime.Time型へ変換します
/////////////////////////////////////////
func getParsedTime(strTime string) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	t, err := time.ParseInLocation(DATETIME_FORMAT, strTime, loc)
	if err != nil {
		return time.Date(1970, 1, 1, 9, 0, 0, 0, loc)
	}

	return t
}

// summary => 現在時刻を示す文字列を取得します
/////////////////////////////////////////
func getNowString() string {
	return time.Now().Format(DATETIME_FORMAT)
}

// summary => 渡された各オクテットごとの数字をIPアドレスへと結合します
/////////////////////////////////////////
func concatIPAddress(a, b, c, d string) string {
	strIP := fmt.Sprintf("%s.%s.%s.%s", a, b, c, d)
	ip := net.ParseIP(strIP)

	if ip == nil {
		return "0.0.0.0"
	} else {
		return strIP
	}
}

