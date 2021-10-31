package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


// summary => 待ち受けるサーバのルータを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("api/v1")

	v1.GET("/init", init)
	v1.GET("/key/:key", continuation)
	v1.POST("/key/:key", next)

	 return router
}

// summary => 最初から始める場合の処理
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func init(c *gin.Context) {
	//途中で処理を中断
	//c.Abort()
}

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

// summary => 固有キーを取得します
// return::string => 一意な固有キー
/////////////////////////////////////////
func getKey() string {
	obj, err := uuid.NewRandom()

	if err == nil {
		return obj.String()
	} else {
		return uuid.Nil.String()
	}
}

