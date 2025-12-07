package context

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/o-ga09/web-ya-hime/pkg/config"
)

type CtxUserKey string
type CtxRequestIDKey string
type CtxDBKey string

const USERID CtxUserKey = "userID"
const REQUESTID CtxRequestIDKey = "requestId"
const DBKEY CtxDBKey = "db"

func GetCtxFromUser(ctx context.Context) string {
	userID, ok := ctx.Value(USERID).(string)
	if !ok {
		return ""
	}
	return userID
}

func SetCtxFromUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, USERID, userID)
}

func SetRequestID(ctx context.Context) context.Context {
	reqID := GetRequestID(ctx)
	if reqID != "" {
		return ctx
	}
	return context.WithValue(ctx, REQUESTID, uuid.NewString())
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(REQUESTID).(string); ok {
		return reqID
	}
	return ""
}

func GetCtxCfg(ctx context.Context) *config.Config {
	cfg, ok := ctx.Value(config.CtxEnvKey).(*config.Config)
	if !ok {
		return &config.Config{}
	}
	return cfg
}

func SetDB(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, DBKEY, db)
}

func GetDB(ctx context.Context) *sql.DB {
	if db, ok := ctx.Value(DBKEY).(*sql.DB); ok {
		return db
	}
	return nil
}
