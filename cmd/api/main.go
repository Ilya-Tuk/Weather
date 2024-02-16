package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/repositories/memory"
	"github.com/Ilya-Tuk/Weather/internal/services"
	"github.com/Ilya-Tuk/Weather/internal/transport/rest"
	"go.uber.org/zap"
)

func main() {
	logCfg := zap.NewDevelopmentConfig()
	logCfg.OutputPaths = []string{"server.log"}
	logCfg.Encoding = "json"

	logger, _ := logCfg.Build()
	defer logger.Sync()
	lg := logger.Sugar()

	cfg := config.Read()

	repo := &memory.Repository{}
	repo.Init()
	service := services.New(repo)

	defer repo.Close()

	server := rest.NewServer(cfg, service, lg)
	lg.Info("Server started:\n	Host:", cfg.ServCfg.ServerHost, "\n	WeatherApi:", cfg.ServCfg.WeatherApiUrl, "\n	Debug Mode:", cfg.ServCfg.DebugMode)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Panicln("Shutdown error:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Panicln(err)
	}
}
