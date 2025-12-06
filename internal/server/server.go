package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
	"github.com/o-ga09/web-ya-hime/pkg/logger"
)

type IServer interface {
	Run(ctx context.Context) error
}

type server struct{}

func NewServer(ctx context.Context) IServer {
	return &server{}
}

func (s *server) Run(ctx context.Context) error {
	cfg := Ctx.GetCtxCfg(ctx)
	engine := http.NewServeMux()
	engine.HandleFunc("/health", healthCheck)

	// サーバーの起動
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
