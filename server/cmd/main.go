package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/darkseear/gophkeeper/proto"
	"github.com/darkseear/gophkeeper/server/internal/api/proto"
	"github.com/darkseear/gophkeeper/server/internal/config"
	"github.com/darkseear/gophkeeper/server/internal/logger"
	"github.com/darkseear/gophkeeper/server/internal/storage"
)

// GRPCServer - структура для gRPC сервера.
type GRPCServer struct {
	Server *grpc.Server
}

// App - основная структура приложения, содержащая серверы, хранилище и конфигурацию.
type App struct {
	GRPCServer *GRPCServer
	Storage    *storage.Store
	Cfg        *config.Config
}

// newApp - инициализирует приложение, настраивает логирование, хранилище и роутер.
func newApp(ctx context.Context) (*App, error) {
	cfg := config.New()
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}
	defer logger.Log.Sync()

	stor, err := storage.NewStore(cfg)
	if err != nil {
		logger.Log.Error("Error store created")
		return nil, err
	}

	// Настройка gRPC сервера
	grpcSrv := grpc.NewServer()
	nss := proto.NewGophkeeperGRPCServer(stor, cfg)
	pb.RegisterGophkeeperServer(grpcSrv, nss)

	return &App{
		GRPCServer: &GRPCServer{
			Server: grpcSrv,
		},
		Storage: stor,
		Cfg:     cfg,
	}, nil
}

// Run - запускает приложение, инициализирует сервер и обрабатывает сигналы завершения.
func (a *App) Run(ctx context.Context) {

	// Запуск gRPC сервера
	go func() {
		listener, err := net.Listen("tcp", a.Cfg.Address)
		if err != nil {
			logger.Log.Error("Error starting gRPC server", zap.Error(err))
			log.Fatalf("Error starting gRPC server: %v", err)
		}
		logger.Log.Info("Starting GRPC server", zap.String("GRPCaddress", a.Cfg.Address))
		err = a.GRPCServer.Server.Serve(listener)
		if err != nil {
			logger.Log.Error("Error starting gRPC server", zap.Error(err))
			log.Fatalf("Error starting gRPC server: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Log.Info("Received shutdown signal, shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := a.Close(shutdownCtx); err != nil {
		logger.Log.Error("Error during shutdown", zap.Error(err))
	} else {
		logger.Log.Info("Server shutdown gracefully")
	}

}

// Close - закрывает приложение, останавливает сервер и освобождает ресурсы.
func (a *App) Close(ctx context.Context) error {
	var errs []error

	// Останавливаем gRPC сервер
	if a.GRPCServer != nil && a.GRPCServer.Server != nil {
		a.GRPCServer.Server.GracefulStop()
		logger.Log.Info("gRPC server stopped")
	}

	// Закрываем storage
	if err := a.Storage.Close(); err != nil {
		logger.Log.Error("Error closing storage", zap.Error(err))
		errs = append(errs, err)
	} else {
		logger.Log.Info("Storage closed")
	}

	// Синхронизируем логгер
	if err := logger.Log.Sync(); err != nil {
		logger.Log.Error("Error syncing logger", zap.Error(err))
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	logger.Log.Info("Application closed successfully")
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	app, err := newApp(ctx)
	if err != nil {
		logger.Log.Error("Error create app", zap.Error(err))
		log.Fatalf("Error app: %v", err)
	}
	app.Run(ctx)
	defer app.Close(ctx)
}
