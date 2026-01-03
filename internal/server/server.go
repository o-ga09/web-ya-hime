package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/o-ga09/web-ya-hime/internal/handler/category"
	"github.com/o-ga09/web-ya-hime/internal/handler/subcategory"
	"github.com/o-ga09/web-ya-hime/internal/handler/summary"
	"github.com/o-ga09/web-ya-hime/internal/handler/user"
	"github.com/o-ga09/web-ya-hime/internal/infra/database/mysql"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
	"github.com/o-ga09/web-ya-hime/pkg/httputil"
	"github.com/o-ga09/web-ya-hime/pkg/logger"
)

type IServer interface {
	Run(ctx context.Context) error
}

type server struct {
	user        user.IUserHandler
	summary     summary.ISummaryHandler
	category    category.ICategoryHandler
	subcategory subcategory.ISubcategoryHandler
}

func NewServer(ctx context.Context) IServer {
	summaryRepo := mysql.NewSummaryRepository()
	userRepo := mysql.NewUserRepository()
	categoryRepo := mysql.NewCategoryRepository()
	subcategoryRepo := mysql.NewSubcategoryRepository()
	return &server{
		user:        user.New(userRepo),
		summary:     summary.New(summaryRepo, subcategoryRepo),
		category:    category.New(categoryRepo),
		subcategory: subcategory.New(subcategoryRepo),
	}
}

func (s *server) Run(ctx context.Context) error {
	cfg := Ctx.GetCtxCfg(ctx)
	engine := http.NewServeMux()

	// ヘルスチェックAPI
	healthCheckHandler := UseMiddleware(ctx, healthCheck)
	DBHealthCheckHandler := UseMiddleware(ctx, DBHealthCheck)

	engine.HandleFunc("/health", healthCheckHandler)
	engine.HandleFunc("/db-health", DBHealthCheckHandler)

	// ユーザーAPI
	userSaveHandler := UseMiddleware(ctx, s.user.Save)
	userListHandler := UseMiddleware(ctx, s.user.List)
	userDetailHandler := UseMiddleware(ctx, s.user.Detail)
	userDeleteHandler := UseMiddleware(ctx, s.user.Delete)

	engine.HandleFunc("POST /users", userSaveHandler)
	engine.HandleFunc("GET /users", userListHandler)
	engine.HandleFunc("GET /users/{id}", userDetailHandler)
	engine.HandleFunc("DELETE /users/{id}", userDeleteHandler)

	// 概要欄取得API
	summarySaveHandler := UseMiddleware(ctx, s.summary.Save)
	summaryListHandler := UseMiddleware(ctx, s.summary.List)
	summaryDetailHandler := UseMiddleware(ctx, s.summary.Detail)
	summaryDeleteHandler := UseMiddleware(ctx, s.summary.Delete)

	engine.HandleFunc("POST /summaries", summarySaveHandler)
	engine.HandleFunc("GET /summaries", summaryListHandler)
	engine.HandleFunc("GET /summaries/{id}", summaryDetailHandler)
	engine.HandleFunc("DELETE /summaries/{id}", summaryDeleteHandler)

	// カテゴリAPI
	categorySaveHandler := UseMiddleware(ctx, s.category.Save)
	categoryListHandler := UseMiddleware(ctx, s.category.List)
	categoryDetailHandler := UseMiddleware(ctx, s.category.Detail)
	categoryDeleteHandler := UseMiddleware(ctx, s.category.Delete)

	engine.HandleFunc("POST /categories", categorySaveHandler)
	engine.HandleFunc("GET /categories", categoryListHandler)
	engine.HandleFunc("GET /categories/{id}", categoryDetailHandler)
	engine.HandleFunc("DELETE /categories/{id}", categoryDeleteHandler)

	// サブカテゴリAPI
	subcategorySaveHandler := UseMiddleware(ctx, s.subcategory.Save)
	subcategoryListHandler := UseMiddleware(ctx, s.subcategory.List)
	subcategoryDetailHandler := UseMiddleware(ctx, s.subcategory.Detail)
	subcategoryDeleteHandler := UseMiddleware(ctx, s.subcategory.Delete)

	engine.HandleFunc("POST /subcategories", subcategorySaveHandler)
	engine.HandleFunc("GET /subcategories", subcategoryListHandler)
	engine.HandleFunc("GET /subcategories/{id}", subcategoryDetailHandler)
	engine.HandleFunc("DELETE /subcategories/{id}", subcategoryDeleteHandler)

	port := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: engine,
	}

	// サーバーの起動
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, fmt.Sprintf("Failed to listen and serve: %v", err))
		}
	}()

	logger.Info(ctx, fmt.Sprintf("Server is running on %s", port))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "graceful shutdown")

	// サーバーのタイムアウト設定
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// サーバーのシャットダウン
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to shutdown server: %v", err))
		return err
	}

	return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	httputil.Response(&w, http.StatusOK, map[string]string{"message": "OK"})
}

func DBHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := Ctx.GetDB(ctx)
	if db == nil {
		http.Error(w, "Database connection not found", http.StatusInternalServerError)
		return
	}

	if err := db.PingContext(ctx); err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	httputil.Response(&w, http.StatusOK, map[string]string{"message": "Database connection is healthy"})
}
