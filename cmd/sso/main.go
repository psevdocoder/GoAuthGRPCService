package main

import (
	"authService/internal/app"
	"authService/internal/config"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log, logfile := SetupLogger(cfg.Env)
	log.Debug("Loaded config", "config", *cfg)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("Application stopped", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()
	log.Info("Application stopped")

	_ = logfile.Close()
}

func SetupLogger(env string) (*slog.Logger, *os.File) {

	logfile, err := os.OpenFile(
		fmt.Sprintf("logs/env_%s.log", env), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	var logger *slog.Logger

	multiWriter := io.MultiWriter(os.Stdout, logfile)

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger, logfile
}
