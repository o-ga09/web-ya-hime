package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
)

const (
	maxRetries      = 5
	retryInterval   = 2 * time.Second
	maxIdleConns    = 10
	maxOpenConns    = 100
	connMaxLifetime = 1 * time.Hour
)

func Connect(ctx context.Context) (*sql.DB, error) {
	var db *sql.DB
	var err error
	env := Ctx.GetCtxCfg(ctx)
	dsn := env.Database_url

	// リトライ処理
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			if i == maxRetries-1 {
				return nil, fmt.Errorf("failed to open database after %d retries: %w", maxRetries, err)
			}
			time.Sleep(retryInterval)
			continue
		}

		// 接続確認
		if err = db.PingContext(ctx); err != nil {
			db.Close()
			if i == maxRetries-1 {
				return nil, fmt.Errorf("failed to ping database after %d retries: %w", maxRetries, err)
			}
			time.Sleep(retryInterval)
			continue
		}

		// 接続成功
		break
	}

	// コネクションプール設定
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	return db, nil
}
