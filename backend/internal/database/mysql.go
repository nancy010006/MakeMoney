package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func SetupDatabase() (*sql.DB, error) {
	var err error
	for i := 0; i < 30; i++ { // 嘗試 30 次，每次等待 2 秒
		db, err = ConnectDatabase()
		if err == nil {
			return db, nil
		}
		log.Printf("嘗試連接數據庫失敗 (嘗試 %d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("設置數據庫失敗，已重試 30 次: %v", err)
}

func ConnectDatabase() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打開數據庫連接失敗: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("連接數據庫失敗: %v", err)
	}

	log.Println("成功連接到數據庫")
	return db, nil
}

func TestDatabaseConnection() error {
	if db == nil {
		return fmt.Errorf("數據庫未初始化")
	}

	// 執行一個簡單的查詢
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("執行測試查詢時出錯: %v", err)
	}

	if result != 1 {
		return fmt.Errorf("測試查詢的意外結果: 得到 %d, 期望 1", result)
	}

	fmt.Println("數據庫連接測試成功")
	return nil
}

func GetDB() *sql.DB {
	return db
}
