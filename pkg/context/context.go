package context

import (
	"context"

	"github.com/google/uuid"
	"github.com/o-ga09/web-ya-hime/pkg/config"
)

type CtxUserKey string
type CtxRequestIDKey string

const USERID CtxUserKey = "userID"
const REQUESTID CtxRequestIDKey = "requestId"

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
