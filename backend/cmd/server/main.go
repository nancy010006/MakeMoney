package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 創建一個默認的 Gin 路由器
	r := gin.Default()

	// 定義一個簡單的路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "歡迎來到賺錢樂！",
		})
	})

	// 添加更多路由
	setupRoutes(r)

	// 啟動服務器
	log.Println("服務器正在啟動，監聽端口 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("啟動服務器失敗: %v", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// 用戶相關路由
	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)

	// 股票相關路由
	r.GET("/stocks", getStocksHandler)
	r.POST("/trade", tradeHandler)

	// 遊戲相關路由
	r.GET("/game/start", startGameHandler)
	r.GET("/game/status", gameStatusHandler)
}

// 以下是處理函數的存根，稍後我們會實現它們
func registerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "註冊功能待實現"})
}

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "登錄功能待實現"})
}

func getStocksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "獲取股票信息功能待實現"})
}

func tradeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "交易功能待實現"})
}

func startGameHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "開始遊戲功能待實現"})
}

func gameStatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "獲取遊戲狀態功能待實現"})
}
