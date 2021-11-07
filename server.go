package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"

	"ip-calc-practice-api/db"
	"ip-calc-practice-api/models"
)

var dbmap *gorp.DbMap

// summary => 待ち受けるサーバのルータを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("api/v1")

	v1.GET("/init", initializeAction)
	//v1.GET("/key/:key", continuation)
	//v1.POST("/key/:key", getNextQuestion)

	dbmap = initDB()

	 return router
}

// summary => 最初から始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func initializeAction(c *gin.Context) {
	if dbmap == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "DBと接続できません",
		})

		return
	}

	repo := db.NewIpRepository(dbmap)
	m, err := repo.GetExpire(context.Background(), "cf50bf5d-c4e1-4088-af5e-ffbd0257e770")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "データの取得に失敗しました",
		})

		return
	}

	c.JSON(http.StatusOK, m)
}

/*
// summary => 続きから始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func continuation(c *gin.Context) {
	//途中で処理を中断
	//c.Abort()
}

// summary => 解答をDBへ格納して、次の問題を返す処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func next(c *gin.Context) {
	//途中で処理を中断
	//c.Abort()
}
*/

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

