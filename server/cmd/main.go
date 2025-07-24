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

// Server - интерфейс для серверов.
type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// GRPCServer - gRPC сервер.
type GRPCServer struct {
	server  *grpc.Server
	address string
}

func NewGRPCServer(address string, stor *storage.Store, cfg *config.Config) *GRPCServer {
	grpcSrv := grpc.NewServer()
	nss := proto.NewGophkeeperGRPCServer(stor, cfg)
	pb.RegisterGophkeeperServer(grpcSrv, nss)
	return &GRPCServer{
		server:  grpcSrv,
		address: address,
	}
}

func (g *GRPCServer) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", g.address)
	if err != nil {
		logger.Log.Error("Error starting gRPC server", zap.Error(err))
		return err
	}
	logger.Log.Info("Starting GRPC server", zap.String("GRPCaddress", g.address))
	go func() {
		if err := g.server.Serve(listener); err != nil {
			logger.Log.Error("gRPC server stopped with error", zap.Error(err))
		}
	}()
	return nil
}

// Stop - остановка сервера .
func (g *GRPCServer) Stop(ctx context.Context) error {
	g.server.GracefulStop()
	logger.Log.Info("gRPC server stopped")
	return nil
}

// StorageComponent - компонент для работы с хранилищем.
type StorageComponent struct {
	store *storage.Store
}

// NewStorageComponent - экземпляр компонента сервера.
func NewStorageComponent(cfg *config.Config) (*StorageComponent, error) {
	stor, err := storage.NewStore(cfg)
	if err != nil {
		logger.Log.Error("Error store created")
		return nil, err
	}
	return &StorageComponent{store: stor}, nil
}

// Close - звкрывает хранилище.
func (s *StorageComponent) Close() error {
	if err := s.store.Close(); err != nil {
		logger.Log.Error("Error closing storage", zap.Error(err))
		return err
	}
	logger.Log.Info("Storage closed")
	return nil
}

// Application - агрегирует компоненты приложения.
type Application struct {
	servers []Server
	storage *StorageComponent
	cfg     *config.Config
}

// NewApplication - экземпляр агригатора приложения.
func NewApplication(ctx context.Context) (*Application, error) {
	cfg := config.New()
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		return nil, err
	}
	defer logger.Log.Sync()

	storageComp, err := NewStorageComponent(cfg)
	if err != nil {
		return nil, err
	}

	grpcServer := NewGRPCServer(cfg.Address, storageComp.store, cfg)

	return &Application{
		servers: []Server{grpcServer},
		storage: storageComp,
		cfg:     cfg,
	}, nil
}

// Run - запуск приложения.
func (a *Application) Run(ctx context.Context) {
	for _, srv := range a.servers {
		if err := srv.Start(ctx); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}

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

// Close - закрывает приложение со всеми его компонентами.
func (a *Application) Close(ctx context.Context) error {
	var errs []error

	for _, srv := range a.servers {
		if err := srv.Stop(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if err := a.storage.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := logger.Log.Sync(); err != nil {
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
	app, err := NewApplication(ctx)
	if err != nil {
		logger.Log.Error("Error create app", zap.Error(err))
		log.Fatalf("Error app: %v", err)
	}
	app.Run(ctx)
	defer app.Close(ctx)
}
