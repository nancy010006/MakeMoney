SET NAMES utf8mb4;
SET CHARSET utf8mb4;

-- 遊戲設置表
CREATE TABLE IF NOT EXISTS game_settings (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '設置ID',
    initial_stock_price DECIMAL(10, 2) NOT NULL COMMENT '初始股票價格',
    total_shares INT NOT NULL COMMENT '總股數',
    initial_cash DECIMAL(15, 2) NOT NULL COMMENT '玩家初始資金',
    matching_interval INT NOT NULL DEFAULT 5 COMMENT '撮合間隔（秒）',
    trading_hours VARCHAR(50) NOT NULL DEFAULT '00:00-24:00' COMMENT '交易時間',
    npc_count INT NOT NULL DEFAULT 3 COMMENT 'NPC數量',
    created_at DATETIME NOT NULL COMMENT '創建時間',
    updated_at DATETIME NOT NULL COMMENT '更新時間'
) COMMENT '遊戲全局設置表';

-- 用戶表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '用戶ID',
    username VARCHAR(50) NOT NULL COMMENT '用戶名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密碼哈希',
    email VARCHAR(100) NOT NULL COMMENT '電子郵件',
    cash_balance DECIMAL(15, 2) NOT NULL COMMENT '現金餘額',
    stock_balance INT NOT NULL DEFAULT 0 COMMENT '股票餘額',
    is_npc BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否為NPC',
    npc_type ENUM('COMPANY', 'TRADER') COMMENT 'NPC類型',
    created_at DATETIME NOT NULL COMMENT '創建時間',
    updated_at DATETIME NOT NULL COMMENT '更新時間'
) COMMENT '用戶信息表';

CREATE UNIQUE INDEX idx_username ON users(username);
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_is_npc ON users(is_npc);

-- NPC策略表
CREATE TABLE IF NOT EXISTS npc_strategies (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '策略ID',
    user_id INT NOT NULL COMMENT '關聯的NPC用戶ID',
    strategy_type ENUM('RANDOM', 'TREND_FOLLOWING', 'CONTRARIAN', 'COMPANY') NOT NULL COMMENT '策略類型',
    parameters JSON COMMENT '策略參數',
    created_at DATETIME NOT NULL COMMENT '創建時間',
    updated_at DATETIME NOT NULL COMMENT '更新時間'
) COMMENT 'NPC交易策略表';

CREATE INDEX idx_npc_user_id ON npc_strategies(user_id);

-- 訂單簿表
CREATE TABLE IF NOT EXISTS order_book (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '訂單ID',
    user_id INT NOT NULL COMMENT '用戶ID',
    order_type ENUM('BUY', 'SELL') NOT NULL COMMENT '訂單類型',
    price DECIMAL(10, 2) NOT NULL COMMENT '價格',
    quantity INT NOT NULL COMMENT '數量',
    status ENUM('ACTIVE', 'FILLED', 'CANCELLED') NOT NULL DEFAULT 'ACTIVE' COMMENT '訂單狀態',
    created_at DATETIME NOT NULL COMMENT '創建時間',
    updated_at DATETIME NOT NULL COMMENT '更新時間'
) COMMENT '訂單簿表';

CREATE INDEX idx_order_user_id ON order_book(user_id);
CREATE INDEX idx_order_type_status ON order_book(order_type, status);
CREATE INDEX idx_order_price ON order_book(price);

-- 交易歷史表
CREATE TABLE IF NOT EXISTS trade_history (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '交易ID',
    buyer_id INT NOT NULL COMMENT '買方ID',
    seller_id INT NOT NULL COMMENT '賣方ID',
    price DECIMAL(10, 2) NOT NULL COMMENT '成交價格',
    quantity INT NOT NULL COMMENT '成交數量',
    trade_time DATETIME NOT NULL COMMENT '成交時間'
) COMMENT '交易歷史表';

CREATE INDEX idx_trade_buyer_id ON trade_history(buyer_id);
CREATE INDEX idx_trade_seller_id ON trade_history(seller_id);
CREATE INDEX idx_trade_time ON trade_history(trade_time);

-- 股票價格歷史表
CREATE TABLE IF NOT EXISTS stock_price_history (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '記錄ID',
    price DECIMAL(10, 2) NOT NULL COMMENT '價格',
    timestamp DATETIME NOT NULL COMMENT '時間戳'
) COMMENT '股票價格歷史表';

CREATE INDEX idx_stock_price_timestamp ON stock_price_history(timestamp);

-- 股利發放記錄表
CREATE TABLE IF NOT EXISTS dividend_history (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '記錄ID',
    amount_per_share DECIMAL(10, 4) NOT NULL COMMENT '每股股利',
    ex_dividend_date DATETIME NOT NULL COMMENT '除息日期',
    payment_date DATETIME NOT NULL COMMENT '發放日期',
    created_at DATETIME NOT NULL COMMENT '創建時間'
) COMMENT '股利發放記錄表';

CREATE INDEX idx_dividend_ex_date ON dividend_history(ex_dividend_date);
CREATE INDEX idx_dividend_payment_date ON dividend_history(payment_date);

-- 遊戲局表
CREATE TABLE IF NOT EXISTS game_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT '遊戲局ID',
    start_time DATETIME NOT NULL COMMENT '開始時間',
    end_time DATETIME COMMENT '結束時間',
    status ENUM('ACTIVE', 'ENDED') NOT NULL DEFAULT 'ACTIVE' COMMENT '遊戲狀態'
) COMMENT '遊戲局表';

CREATE INDEX idx_game_session_status ON game_sessions(status);

-- 封存表

-- 封存用戶表
CREATE TABLE IF NOT EXISTS archived_users (
    id INT NOT NULL COMMENT '用戶ID',
    game_session_id INT NOT NULL COMMENT '遊戲局ID',
    username VARCHAR(50) NOT NULL COMMENT '用戶名',
    email VARCHAR(100) NOT NULL COMMENT '電子郵件',
    final_cash_balance DECIMAL(15, 2) NOT NULL COMMENT '最終現金餘額',
    final_stock_balance INT NOT NULL COMMENT '最終股票餘額',
    is_npc BOOLEAN NOT NULL COMMENT '是否為NPC',
    npc_type ENUM('COMPANY', 'TRADER') COMMENT 'NPC類型',
    created_at DATETIME NOT NULL COMMENT '創建時間',
    PRIMARY KEY (id, game_session_id)
) COMMENT '封存用戶表';

CREATE INDEX idx_archived_users_game_session ON archived_users(game_session_id);

-- 封存NPC策略表
CREATE TABLE IF NOT EXISTS archived_npc_strategies (
    id INT NOT NULL COMMENT '策略ID',
    game_session_id INT NOT NULL COMMENT '遊戲局ID',
    user_id INT NOT NULL COMMENT '關聯的NPC用戶ID',
    strategy_type ENUM('RANDOM', 'TREND_FOLLOWING', 'CONTRARIAN', 'COMPANY') NOT NULL COMMENT '策略類型',
    parameters JSON COMMENT '策略參數',
    PRIMARY KEY (id, game_session_id)
) COMMENT '封存NPC策略表';

CREATE INDEX idx_archived_npc_strategies_game_session ON archived_npc_strategies(game_session_id);

-- 封存訂單簿表
CREATE TABLE IF NOT EXISTS archived_order_book (
    id INT NOT NULL COMMENT '訂單ID',
    game_session_id INT NOT NULL COMMENT '遊戲局ID',
    user_id INT NOT NULL COMMENT '用戶ID',
    order_type ENUM('BUY', 'SELL') NOT NULL COMMENT '訂單類型',
    price DECIMAL(10, 2) NOT NULL COMMENT '價格',
    quantity INT NOT NULL COMMENT '數量',
    status ENUM('ACTIVE', 'FILLED', 'CANCELLED') NOT NULL COMMENT '訂單狀態',
    created_at DATETIME NOT NULL COMMENT '創建時間',
    updated_at DATETIME NOT NULL COMMENT '更新時間',
    PRIMARY KEY (id, game_session_id)
) COMMENT '封存訂單簿表';

CREATE INDEX idx_archived_order_book_game_session ON archived_order_book(game_session_id);

-- 封存交易歷史表
CREATE TABLE IF NOT EXISTS archived_trade_history (
    id INT NOT NULL COMMENT '交易ID',
    game_session_id INT NOT NULL COMMENT '遊戲局ID',
    buyer_id INT NOT NULL COMMENT '買方ID',
    seller_id INT NOT NULL COMMENT '賣方ID',
    price DECIMAL(10, 2) NOT NULL COMMENT '成交價格',
    quantity INT NOT NULL COMMENT '成交數量',
    trade_time DATETIME NOT NULL COMMENT '成交時間',
    PRIMARY KEY (id, game_session_id)
) COMMENT '封存交易歷史表';

CREATE INDEX idx_archived_trade_history_game_session ON archived_trade_history(game_session_id);

-- 封存股票價格歷史表
CREATE TABLE IF NOT EXISTS archived_stock_price_history (
    id INT NOT NULL COMMENT '記錄ID',
    game_session_id INT NOT NULL COMMENT '遊戲局ID',
    price DECIMAL(10, 2) NOT NULL COMMENT '價格',
    timestamp DATETIME NOT NULL COMMENT '時間戳',
    PRIMARY KEY (id, game_session_id)
) COMMENT '封存股票價格歷史表';

CREATE INDEX idx_archived_stock_price_game_session ON archived_stock_price_history(game_session_id);

-- 封存股利發放記錄表
CREATE TABLE IF NOT EXISTS archived_dividend_history (
    id INT NOT NULL COMMENT '記錄ID',
    game_session_id INT NOT NULL COMMENT '遊戲局ID',
    amount_per_share DECIMAL(10, 4) NOT NULL COMMENT '每股股利',
    ex_dividend_date DATETIME NOT NULL COMMENT '除息日期',
    payment_date DATETIME NOT NULL COMMENT '發放日期',
    created_at DATETIME NOT NULL COMMENT '創建時間',
    PRIMARY KEY (id, game_session_id)
) COMMENT '封存股利發放記錄表';

CREATE INDEX idx_archived_dividend_history_game_session ON archived_dividend_history(game_session_id);

-- 插入初始遊戲設置
INSERT INTO game_settings (initial_stock_price, total_shares, initial_cash, matching_interval, trading_hours, npc_count, created_at, updated_at)
VALUES (100.00, 1000000, 10000.00, 5, '00:00-24:00', 3, NOW(), NOW());

-- 插入NPC用戶
INSERT INTO users (username, password_hash, email, cash_balance, stock_balance, is_npc, npc_type, created_at, updated_at) VALUES
('NPC_Company', 'hashed_password', 'npc_company@example.com', 0, 1000000, TRUE, 'COMPANY', NOW(), NOW()),
('NPC_Trader1', 'hashed_password', 'npc_trader1@example.com', 10000, 0, TRUE, 'TRADER', NOW(), NOW()),
('NPC_Trader2', 'hashed_password', 'npc_trader2@example.com', 10000, 0, TRUE, 'TRADER', NOW(), NOW()),
('NPC_Trader3', 'hashed_password', 'npc_trader3@example.com', 10000, 0, TRUE, 'TRADER', NOW(), NOW());

-- 插入NPC策略
INSERT INTO npc_strategies (user_id, strategy_type, parameters, created_at, updated_at) VALUES
((SELECT id FROM users WHERE username = 'NPC_Company'), 'COMPANY', '{"dividend_rate": 0.05, "dividend_interval": "3 months"}', NOW(), NOW()),
((SELECT id FROM users WHERE username = 'NPC_Trader1'), 'RANDOM', '{"min_order": 10, "max_order": 100}', NOW(), NOW()),
((SELECT id FROM users WHERE username = 'NPC_Trader2'), 'TREND_FOLLOWING', '{"lookback_period": 5}', NOW(), NOW()),
((SELECT id FROM users WHERE username = 'NPC_Trader3'), 'CONTRARIAN', '{"threshold": 0.05}', NOW(), NOW());

-- 插入初始遊戲局
INSERT INTO game_sessions (start_time, status) VALUES (NOW(), 'ACTIVE');