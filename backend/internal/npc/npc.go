package npc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"
)

type NPC struct {
	ID     string
	db     *sql.DB
	cancel context.CancelFunc
}

func New(id string, db *sql.DB) *NPC {
	return &NPC{
		ID: id,
		db: db,
	}
}

func (n *NPC) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	n.cancel = cancel

	go n.run(ctx)
}

func (n *NPC) Stop() {
	if n.cancel != nil {
		n.cancel()
	}
}

func (n *NPC) run(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := n.generateOrder()
			if err != nil {
				log.Printf("NPC %s 生成訂單時出錯: %v", n.ID, err)
			}
		}
	}
}

func (n *NPC) generateOrder() error {
	// 獲取最新成交價
	var lastPrice float64
	err := n.db.QueryRow("SELECT price FROM stock_price_history ORDER BY timestamp DESC LIMIT 1").Scan(&lastPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			err = n.db.QueryRow("SELECT initial_stock_price FROM game_settings ORDER BY id DESC LIMIT 1").Scan(&lastPrice)
			if err != nil {
				return fmt.Errorf("獲取初始股票價格失敗: %v", err)
			}
		} else {
			return fmt.Errorf("獲取最新股票價格失敗: %v", err)
		}
	}

	// 生成訂單
	orderType := "SELL"
	if rand.Float32() < 0.5 {
		orderType = "BUY"
	}

	price := lastPrice
	if orderType == "SELL" {
		price *= 1.05 // 賣出價比市場價高5%
	} else {
		price *= 0.95 // 買入價比市場價低5%
	}

	quantity := rand.Intn(40001) + 10000 // 10000 到 50000 之間的隨機數

	// 將訂單插入數據庫
	_, err = n.db.Exec(
		"INSERT INTO order_book (user_id, order_type, price, quantity, status, created_at, updated_at) VALUES (?, ?, ?, ?, 'ACTIVE', NOW(), NOW())",
		n.ID, orderType, price, quantity,
	)
	if err != nil {
		return fmt.Errorf("插入訂單失敗: %v", err)
	}

	log.Printf("NPC %s 生成訂單: 類型=%s, 價格=%.2f, 數量=%d", n.ID, orderType, price, quantity)
	return nil
}
