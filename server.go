package main

import (
	_ "context"
	"database/sql"
	"fmt"
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

var dbmap *gorp.DbMap

// summary => 待ち受けるサーバのルータを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("api/v1")

	v1.GET("/init/:total", initializeAction)
	//v1.GET("/resume/:key", resumeAnswer)
	v1.POST("/next/:key", getNextQuestion)
	//v1.DELETE("clean", deleteExpiredData)

	dbmap = initDB()

	 return router
}

// summary => 最初から始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func initializeAction(c *gin.Context) {
	total_prm := c.Param("total")
	total := getTotalValue(total_prm)

	_, err := c.Cookie(COOKIE_NAME)
	if err == nil {
		c.SetCookie(COOKIE_NAME, "", -1, "/", "", true, true)
	}

	ipint, bits := getQuestion()
	id := getID()
	ip := uint2ip(ipint)

	mid := models.MstrID{
		Id:     id,
		Total:  total,
		Expire: time.Now().Format(DATETIME_FORMAT),
	}

	tq := models.TranQuestion{
		Id:        id,
		Number:    1,
		Source:    ip.String(),
		CIDRbits:  bits,
		IsCIDR:    int(ipint % 2),
		CorNwAddr: getNetworkAddress(ip, bits).String(),
		AnsNwAddr: "0.0.0.0",
		CorBcAddr: getBroadcastAddress(ip, bits).String(),
		AnsBcAddr: "0.0.0.0",
		Elapsed:   0,
	}

	if dbmap == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error"  : "E001",
			"message": "DBと接続できません",
		})

		c.Abort()
	}

	err_fail_data_operation := func() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error"  : "E002",
			"message": "データの操作に失敗しました",
		})

		c.Abort()
	}

	tx, err := dbmap.Begin()

	if err != nil {
		err_fail_data_operation()
	}

	if err := tx.Insert(&mid); err != nil {
		err_fail_data_operation()
	}

	if err := tx.Insert(&tq); err != nil {
		err_fail_data_operation()
	}

	if err := tx.Commit(); err != nil {
		err_fail_data_operation()
	}

	c.SetCookie(COOKIE_NAME, id, 86400, "/", "", true, true)

	res := models.QuestionSet{
		Id:         id,
		Number:     tq.Number,
		Source:     tq.Source,
		CIDRbits:   tq.CIDRbits,
		IsCIDR:     tq.IsCIDR,
		SubnetMask: "",
	}

	if tq.IsCIDR == 1 {
		res.SubnetMask = getSubnetMask(tq.CIDRbits)
	}

	c.JSON(http.StatusOK, res)
}

/*
// summary => 続きから始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func resumeAnswer(c *gin.Context) {
	//途中で処理を中断
	//c.Abort()
	repo := db.NewIpRepository(dbmap)
	m, err := repo.GetExpire(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "データの取得に失敗しました",
		})

		return
	}
}
*/

// summary => 解答をDBへ格納して、次の問題 or 結果を返す処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func getNextQuestion(c *gin.Context) {
}

func initDB() *gorp.DbMap {
	driver, dsn, err := db.GetDataSourceName()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var init_dbmap *gorp.DbMap

	switch driver {
	case "mysql":
		db, _ := sql.Open(driver, dsn)
		dial  := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}

		init_dbmap = &gorp.DbMap{Db: db, Dialect: dial}
		models.MapStructsToTables(init_dbmap)
	}

	return init_dbmap
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
// return::*int => 出題数
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

func getParsedTime(strTime string) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	t, err := time.ParseInLocation(DATETIME_FORMAT, strTime, loc)
	if err != nil {
		return time.Date(1970, 1, 1, 9, 0, 0, 0, loc)
	}

	return t
}

