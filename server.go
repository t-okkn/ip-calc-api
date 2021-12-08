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
	v1 := router.Group("v1")

	v1.GET("/hasCookie", checkCookieStatus)
	v1.GET("/init/:total", initializeAction)
	v1.POST("/next/:id/:number", getNextQuestion)
	v1.GET("/resume", resumeAnswer)
	v1.PUT("/:id/:number", updateQuestion)
	v1.GET("/:id", getRegisteredQuestion)
	v1.GET("/:id/:number", getRegisteredQuestion)
	v1.DELETE("/clean", deleteExpiredData)

	repo = initDB()

	return router
}

// summary => Cookieの確認処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func checkCookieStatus(c *gin.Context) {
	_, err := c.Cookie(COOKIE_NAME)

	if err != nil {
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}
}

// summary => 最初から始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func initializeAction(c *gin.Context) {
	total_prm := c.Param("total")
	total := getTotalValue(total_prm)

	id := getID()
	tq := generateNewQuestion(id, 1)

	tid := models.TranID{
		Id:     id,
		Total:  total,
		IsEnd:  0,
		Expire: time.Now().AddDate(0, 0, 1).Format(DATETIME_FORMAT),
	}

	if repo == nil {
		c.JSON(http.StatusServiceUnavailable, errCannotConnectDB)
		c.Abort()
		return
	}

	if err := repo.InsertFirstData(tid, tq); err != nil {
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errInvalidRequestedData)
		c.Abort()
		return
	}

	readyForUpdateQuestion(req, &tq)

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
		tid, err := repo.GetID(id_prm)

		if err != nil {
			c.JSON(http.StatusBadRequest, errFailedGetData)
			c.Abort()
			return
		}

		tid.IsEnd = 1

		if err := repo.UpdateID(tid); err != nil {
			c.JSON(http.StatusServiceUnavailable, errFailedOperateData)
			c.Abort()
			return
		}

		results, err := repo.GetResults(id_prm)

		if err != nil {
			c.JSON(http.StatusBadRequest, errFailedGetData)
			c.Abort()
			return
		}

		sc := getSummaryCollection(tid, results)
		sc.IsEnd = true

		c.SetCookie(COOKIE_NAME, "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, sc)
	}
}

// summary => 続きから始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func resumeAnswer(c *gin.Context) {
	id, err := c.Cookie(COOKIE_NAME)

	if err != nil {
		c.SetCookie(COOKIE_NAME, "", -1, "/", "", false, true)

		c.JSON(http.StatusBadRequest, errInvalidCookie)
		c.Abort()
		return
	}

	if repo == nil {
		c.JSON(http.StatusServiceUnavailable, errCannotConnectDB)
		c.Abort()
		return
	}

	tid, err := repo.GetID(id)

	if err != nil {
		c.SetCookie(COOKIE_NAME, "", -1, "/", "", false, true)

		c.JSON(http.StatusBadRequest, errRequestedCookieIsNotExist)
		c.Abort()
		return
	}

	if tid.IsEnd == 1 || getParsedTime(tid.Expire).Before(time.Now()) {
		c.SetCookie(COOKIE_NAME, "", -1, "/", "", false, true)

		c.JSON(http.StatusBadRequest, errInvalidRequestedCookie)
		c.Abort()
		return
	}

	tid.Expire = time.Now().AddDate(0, 0, 1).Format(DATETIME_FORMAT)

	if err := repo.UpdateID(tid); err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedOperateData)
		c.Abort()
		return
	}

	qnum, err := repo.CheckNow(id)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedGetData)
		c.Abort()
		return
	}

	tq, err := repo.GetQuestion(id, qnum.Now)

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedGetData)
		c.Abort()
		return
	}

	if qnum.Now > 1 {
		old_tq, err := repo.GetQuestion(id, qnum.Now-1)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, errFailedGetData)
			c.Abort()
			return
		}

		tq.Elapsed = old_tq.Elapsed
	}

	c.SetCookie(COOKIE_NAME, id, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, getResumeSet(tq))
}

// summary => DBの解答情報を更新する処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func updateQuestion(c *gin.Context) {
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

	tid, err := repo.GetID(id_prm)

	if err != nil {
		c.JSON(http.StatusBadRequest, errFailedGetData)
		c.Abort()
		return
	}

	if tid.IsEnd == 1 {
		c.JSON(http.StatusNotAcceptable, errTheQuestionIsTerminated)
		c.Abort()
		return
	}

	check_qnum, err := repo.CheckNow(id_prm)

	if err != nil {
		c.JSON(http.StatusBadRequest, errFailedGetData)
		c.Abort()
		return
	}

	if number > check_qnum.Now {
		c.JSON(http.StatusBadRequest, errTheQuestionIsNotExist)
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errInvalidRequestedData)
		c.Abort()
		return
	}

	readyForUpdateQuestion(req, &tq)

	if err := repo.UpdateQuestion(tq); err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedOperateData)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, sucUpdateDone)
}

// summary => DBに登録されている問題情報を取得する処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func getRegisteredQuestion(c *gin.Context) {
	id_prm  := c.Param("id")
	num_prm := c.Param("number")
	secret  := c.DefaultQuery("secret", "")

	if repo == nil {
		c.JSON(http.StatusServiceUnavailable, errCannotConnectDB)
		c.Abort()
		return
	}

	tid, err := repo.GetID(id_prm)

	if err != nil {
		c.JSON(http.StatusBadRequest, errFailedGetData)
		c.Abort()
		return
	}

	if num_prm != "" {
		number, err := strconv.Atoi(num_prm)

		if err != nil || (number < 1 || number > LIMIT_TOTAL) {
			c.JSON(http.StatusBadRequest, errInvalidRequestedURL)
			c.Abort()
			return
		}

		check_qnum, err := repo.CheckNow(id_prm)

		if err != nil {
			c.JSON(http.StatusBadRequest, errFailedGetData)
			c.Abort()
			return
		}

		if number > check_qnum.Now {
			c.JSON(http.StatusBadRequest, errTheQuestionIsNotExist)
			c.Abort()
			return
		}

		tq, err := repo.GetQuestion(id_prm, number)

		if err != nil {
			c.JSON(http.StatusBadRequest, errFailedGetData)
			c.Abort()
			return
		}

		switch secret {
		case "raw":
			c.JSON(http.StatusOK, getOneRawData(tid, tq))

		case "summary":
			c.JSON(http.StatusOK, getOneSummary(tid, tq))

		default:
			c.JSON(http.StatusOK, getOneResult(tq))
		}

	} else {
		results, err := repo.GetResults(id_prm)

		if err != nil {
			c.JSON(http.StatusBadRequest, errFailedGetData)
			c.Abort()
			return
		}

		switch secret {
		case "raw":
			c.JSON(http.StatusOK, getRawData(tid, results))

		case "summary":
			c.JSON(http.StatusOK, getSummaryCollection(tid, results))

		default:
			c.JSON(http.StatusOK, getResultCollection(results))
		}
	}
}

// summary => DBに登録されている期限切れデータを削除する処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func deleteExpiredData(c *gin.Context) {
	if repo == nil {
		c.JSON(http.StatusServiceUnavailable, errCannotConnectDB)
		c.Abort()
		return
	}

	expiered, err := repo.GetExpired()

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedGetData)
		c.Abort()
		return
	}

	if err := repo.DeleteExpiredData(expiered); err != nil {
		c.JSON(http.StatusServiceUnavailable, errFailedOperateData)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, sucDeleteDone)
}


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

func readyForUpdateQuestion(as models.AnswerSet, tq *models.TranQuestion) {
	str_nwaddr := fmt.Sprintf("%s.%s.%s.%s",
	                          as.NwAddr1st,
	                          as.NwAddr2nd,
	                          as.NwAddr3rd,
	                          as.NwAddr4th)
	str_bcaddr := fmt.Sprintf("%s.%s.%s.%s",
	                          as.BcAddr1st,
	                          as.BcAddr2nd,
	                          as.BcAddr3rd,
	                          as.BcAddr4th)

	nw := net.ParseIP(str_nwaddr)
	bc := net.ParseIP(str_bcaddr)

	var nwaddr string
	if nw == nil {
		nwaddr = "0.0.0.0"
	} else {
		nwaddr = str_nwaddr
	}

	var bcaddr string
	if bc == nil {
		bcaddr = "0.0.0.0"
	} else {
		bcaddr = str_bcaddr
	}

	elapsed, _ := strconv.Atoi(as.Elapsed)

	tq.AnsNwAddr = nwaddr
	tq.AnsBcAddr = bcaddr
	tq.Elapsed   = elapsed
	tq.Updated   = getNowString()
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
func getResumeSet(tq models.TranQuestion) models.ResumeSet {
	rs := models.ResumeSet{
		Id        : tq.Id,
		Number    : tq.Number,
		Source    : tq.Source,
		CIDRbits  : tq.CIDRbits,
		SubnetMask: "",
		Elapsed   : tq.Elapsed,
	}

	if tq.IsCIDR == 1 {
		rs.CIDRbits   = -1
		rs.SubnetMask = getSubnetMask(tq.CIDRbits)
	}

	return rs
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getOneRawData(tid models.TranID, tq models.TranQuestion) models.RawData {
	return models.RawData{
		Id      : tid,
		Question: []models.TranQuestion{tq},
	}
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getRawData(tid models.TranID, tqList []models.TranQuestion) models.RawData {
	return models.RawData{
		Id      : tid,
		Question: tqList,
	}
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getOneSummary(tid models.TranID, tq models.TranQuestion) models.SummaryCollection {
	ss := models.SummarySet{
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
		ss.CIDRbits   = -1
		ss.SubnetMask = getSubnetMask(tq.CIDRbits)
	}

	if tq.Number > 1 {
		if tq.Elapsed != 0 {
			forward_elapsed := 0

			if q, err := repo.GetQuestion(tq.Id, tq.Number-1); err == nil {
				forward_elapsed = q.Elapsed
			}

			if tq.Elapsed - forward_elapsed >= 0 {
				ss.AnswerdTime = tq.Elapsed - forward_elapsed
			}
		}

	} else if tq.Number == 1 {
		ss.AnswerdTime = tq.Elapsed
	}

	return models.SummaryCollection{
		Id     : tq.Id,
		IsEnd  : tid.IsEnd == 1,
		Summary: []models.SummarySet{ss},
	}
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getSummaryCollection(
	tid models.TranID,
	tqList []models.TranQuestion,
) models.SummaryCollection {

	s   := make([]models.SummarySet, len(tqList))
	res := models.SummaryCollection{
		Id     : "",
		IsEnd  : tid.IsEnd == 1,
		Summary: s,
	}

	for i, tq := range tqList {
		ss := models.SummarySet{
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
			ss.CIDRbits   = -1
			ss.SubnetMask = getSubnetMask(tq.CIDRbits)
		}

		if i == 0 {
			res.Id = tq.Id
			ss.AnswerdTime = tq.Elapsed

		} else {
			if tq.Elapsed != 0 {
				if tq.Elapsed - tqList[i-1].Elapsed >= 0 {
					ss.AnswerdTime = tq.Elapsed - tqList[i-1].Elapsed
				}
			}
		}

		res.Summary[i] = ss
	}

	return res
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getOneResult(tq models.TranQuestion) models.ResultCollection {
	rs := models.ResultSet{
		Number     : tq.Number,
		Source     : tq.Source,
		CIDRbits   : tq.CIDRbits,
		SubnetMask : "",
		AnsNwAddr  : tq.AnsNwAddr,
		AnsBcAddr  : tq.AnsBcAddr,
	}

	if tq.IsCIDR == 1 {
		rs.CIDRbits   = -1
		rs.SubnetMask = getSubnetMask(tq.CIDRbits)
	}

	return models.ResultCollection{
		Id    : tq.Id,
		Result: []models.ResultSet{rs},
	}
}

// summary => JSONとして返却するための構造体へデータを詰め替えます
/////////////////////////////////////////
func getResultCollection(tqList []models.TranQuestion) models.ResultCollection {
	r   := make([]models.ResultSet, len(tqList))
	res := models.ResultCollection{
		Id    : "",
		Result: r,
	}

	for i, tq := range tqList {
		if i == 0 {
			res.Id = tq.Id
		}

		rs := models.ResultSet{
			Number     : tq.Number,
			Source     : tq.Source,
			CIDRbits   : tq.CIDRbits,
			SubnetMask : "",
			AnsNwAddr  : tq.AnsNwAddr,
			AnsBcAddr  : tq.AnsBcAddr,
		}

		if tq.IsCIDR == 1 {
			rs.CIDRbits   = -1
			rs.SubnetMask = getSubnetMask(tq.CIDRbits)
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
