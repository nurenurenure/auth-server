package app

import (
	grpcapp "gRPC/api/internal/app/grpc"
	"gRPC/api/internal/services/auth"
	"gRPC/api/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	Storage *postgres.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string, // это DSN строка для PostgreSQL
	tokenTTL time.Duration,
) *App {
	// 1. Подключаемся к PostgreSQL
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// 2. Создаем сервис авторизации
	authService := auth.New(
		log,
		storage, // userSaver
		storage, // userProvider
		storage, // appProvider
		tokenTTL,
	)

	// 3. Создаем gRPC сервер с сервисом авторизации
	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
		Storage: storage,
	}
}
