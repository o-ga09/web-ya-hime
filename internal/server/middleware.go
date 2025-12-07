package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/o-ga09/web-ya-hime/internal/infra/database/mysql"
	"github.com/o-ga09/web-ya-hime/pkg/constant"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
	"github.com/o-ga09/web-ya-hime/pkg/logger"
	"github.com/o-ga09/web-ya-hime/pkg/uuid"
)

type RequestId string

const RequestIdKey RequestId = "requestId"

// AddIDはリクエスト毎にIDを付与するmiddlewareです。
func AddID(ctx context.Context, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IDを生成してcontextに保存
		id := uuid.GenerateID()
		ctx := context.WithValue(ctx, RequestIdKey, id)
		// 次のハンドラーに渡す
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithTimeoutはIDを追加するmiddlewareです。
func WithTimeout(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context() == nil {
			r = r.WithContext(context.Background())
		}

		// タイムアウトを設定
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel() // 処理が終了したらキャンセルする

		// 次のハンドラーを実行し、タイムアウトが発生した場合はエラーメッセージを出力
		done := make(chan struct{})
		go func() {
			defer close(done)
			next.ServeHTTP(w, r)
		}()
		select {
		case <-done:
			// ハンドラーが正常に終了した場合は何もしない
			return
		case <-ctx.Done():
			http.Error(w, "Timeout", http.StatusRequestTimeout)
		}
	})
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(RequestIdKey).(string)
}

type RequestInfo struct {
	ContentsLength int64
	Path           string
	SourceIP       string
	Query          string
	UserAgent      string
	Errors         string
	Elapsed        time.Duration
}

func (r *RequestInfo) LogValue() interface{} { // Assuming slog expects an interface{}
	return map[string]interface{}{
		"contents_length": r.ContentsLength,
		"path":            r.Path,
		"sourceIP":        r.SourceIP,
		"query":           r.Query,
		"user_agent":      r.UserAgent,
		"errors":          r.Errors,
		"elapsed":         r.Elapsed.String(),
	}
}

func RequestLogger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Log(r.Context(), constant.SeverityInfo, "処理開始", "requestId", GetRequestID(r.Context()))
		start := time.Now()

		next.ServeHTTP(w, r)

		req := RequestInfo{
			ContentsLength: r.ContentLength,
			Path:           r.RequestURI,
			SourceIP:       r.RemoteAddr,
			Query:          r.URL.RawQuery,
			UserAgent:      r.UserAgent(),
			Errors:         "errors",
			Elapsed:        time.Since(start),
		}

		slog.Log(r.Context(), constant.SeverityInfo, "処理終了", "Request", req.LogValue(), "requestId", GetRequestID(r.Context())) // Adjust logging context as needed
	})
}

// traceId , spanId 追加
type traceHandler struct {
	slog.Handler
	projectID string
}

// traceHandler 実装
func (h *traceHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	trace := fmt.Sprintf("projects/%s/traces/%s", h.projectID, uuid.GenerateID())
	r.AddAttrs(slog.String("logging.googleapis.com/trace", trace),
		slog.String("logging.googleapis.com/spanId", uuid.GenerateID()))

	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttr(attrs []slog.Attr) slog.Handler {
	return &traceHandler{h.Handler.WithAttrs(attrs), h.projectID}
}

func (h *traceHandler) WithGroup(g string) slog.Handler {
	return h.Handler.WithGroup(g)
}

// logger 生成関数
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		replacer := func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.MessageKey {
				a.Key = "message"
			}

			if a.Key == slog.LevelKey {
				a.Key = "severity"
				a.Value = slog.StringValue(a.Value.Any().(slog.Level).String())
			}

			if a.Key == slog.SourceKey {
				a.Key = "logging.googleapis.com/sourceLocation"
			}

			return a
		}
		env := Ctx.GetCtxCfg(r.Context())
		projectID := env.ProjectID
		h := traceHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replacer}), projectID}
		newh := h.WithAttr([]slog.Attr{
			slog.Group("logging.googleapis.com/labels",
				slog.String("app", "MH-API"),
				slog.String("env", env.Env),
			),
		})
		logger := slog.New(newh)
		slog.SetDefault(logger)
		next.ServeHTTP(w, r)
	})
}

func Cors(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Csrf(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CSRF対策の処理をここに実装
		next.ServeHTTP(w, r)
	})
}

// dbMiddleware はリクエストごとにDBセッションをcontextに設定するミドルウェア
func DBSetUp(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		db, err := mysql.Connect(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to connect to database", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// DBセッションをcontextに設定
		ctx = Ctx.SetDB(ctx, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UseMiddleware(ctx context.Context, handler http.HandlerFunc) http.HandlerFunc {
	handler = WithTimeout(handler)
	handler = DBSetUp(handler)
	handler = RequestLogger(handler)
	handler = Csrf(handler)
	handler = Cors(handler)
	handler = AddID(ctx, handler)
	handler = Logger(handler)

	return handler
}
