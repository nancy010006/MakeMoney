package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase() (*sql.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %v", err)
	}

	// 設置連接池
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 測試連接
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	return db, nil
}

func TestDatabaseConnection() error {
	db, err := ConnectDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// 執行一個簡單的查詢
	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("執行測試查詢時出錯: %v", err)
	}

	if result != 1 {
		return fmt.Errorf("測試查詢的意外結果: 得到 %d, 期望 1", result)
	}

	fmt.Println("數據庫連接測試成功") // 這行會輸出成功消息
	return nil
}
