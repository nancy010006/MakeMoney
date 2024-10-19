package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nancy010006/MakeMoney/internal/database"
	"github.com/nancy010006/MakeMoney/internal/npcmanager"
	"github.com/nancy010006/MakeMoney/internal/services"
)

type App struct {
	db                *sql.DB
	tradingSimulation *services.TradingSimulation
	npcMgr            *npcmanager.NPCManager
	server            *http.Server
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := &App{}
	app.db = setupDatabase()
	defer app.db.Close()

	app.tradingSimulation = setupTradingSimulation(ctx, app.db)
	app.npcMgr = npcmanager.New(app.db)

	router := setupRouter(app.npcMgr)
	app.server = setupServer(router)

	go startServer(app.server)

	waitForShutdown(ctx, app)
}

func setupDatabase() *sql.DB {
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatalf("設置數據庫失敗: %v", err)
	}
	return db
}

func setupTradingSimulation(ctx context.Context, db *sql.DB) *services.TradingSimulation {
	tradingSimulation := services.NewTradingSimulation(db)
	go func() {
		if err := tradingSimulation.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("交易模擬服務停止: %v", err)
		}
	}()
	return tradingSimulation
}

func setupRouter(npcMgr *npcmanager.NPCManager) *gin.Engine {
	r := gin.Default()
	setupRoutes(r, npcMgr)
	return r
}

func setupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}

func startServer(srv *http.Server) {
	log.Println("服務器正在啟動，監聽端口 8080...")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("啟動服務器失敗: %v", err)
	}
}

func waitForShutdown(ctx context.Context, app *App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在關閉服務器...")

	app.npcMgr.StopAll()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := app.server.Shutdown(ctx); err != nil {
		log.Fatalf("服務器被強制關閉: %v", err)
	}

	log.Println("服務器已退出")
}

func setupRoutes(r *gin.Engine, npcMgr *npcmanager.NPCManager) {
	// 定義一個簡單的路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "歡迎來到賺錢樂！",
		})
	})

	// 用戶相關路由
	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)

	// 股票相關路由
	r.GET("/stocks", getStocksHandler)
	r.POST("/trade", tradeHandler)

	// 遊戲相關路由
	r.GET("/game/start", startGameHandler)
	r.GET("/game/status", gameStatusHandler)

	// NPC 管理路由
	r.POST("/npc/start", npcMgr.StartNPCHandler)
	r.POST("/npc/stop", npcMgr.StopNPCHandler)
	r.GET("/npc/list", npcMgr.ListNPCsHandler)
}

// 以下是處理函數的存根，保持不變
func registerHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "註冊功能待實現"})
}

func loginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "登錄功能待實現"})
}

func getStocksHandler(c *gin.Context) {
	// 測試數據庫連接
	err := database.TestDatabaseConnection()
	if err != nil {
		log.Printf("數據庫連接測試失敗: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "數據庫連接失敗"})
		return
	}
	log.Println("數據庫連接測試成功")
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
