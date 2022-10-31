package main

import (
	"fmt"
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	"github.com/eneskzlcn/currency-conversion-service/internal/config"
	"github.com/eneskzlcn/currency-conversion-service/internal/conversion"
	"github.com/eneskzlcn/currency-conversion-service/internal/exchange"
	"github.com/eneskzlcn/currency-conversion-service/internal/wallet"
	"github.com/eneskzlcn/currency-conversion-service/logger"
	"github.com/eneskzlcn/currency-conversion-service/postgres"
	"github.com/eneskzlcn/currency-conversion-service/server"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	env, exists := os.LookupEnv("DEPLOYMENT_ENV")
	if !exists {
		env = "local"
	}
	configs, err := config.LoadConfig(".dev", env, "yaml")
	if err != nil {
		return err
	}

	logger, err := logger.NewZapLoggerForEnv(env, 0)
	if err != nil {
		return err
	}

	db, err := postgres.New(configs.Db)
	if err != nil {
		return err
	}

	authRepository := auth.NewRepository(db)
	authService := auth.NewService(configs.Jwt, authRepository)
	authHandler := auth.NewHandler(authService)
	authGuard := auth.NewGuard(authService, logger)

	walletRepository := wallet.NewRepository(db)
	walletService := wallet.NewService(walletRepository)
	walletHandler := wallet.NewHandler(walletService, authGuard)

	exchangeRepository := exchange.NewRepository(db)
	exchangeService := exchange.NewService(exchangeRepository)
	exchangeHandler := exchange.NewHandler(exchangeService, authGuard)

	conversionService := conversion.NewService(walletService)
	conversionHandler := conversion.NewHandler(conversionService, authGuard)

	server := server.New([]server.Handler{
		authHandler,
		walletHandler,
		exchangeHandler,
		conversionHandler,
	}, configs.Server, logger)

	return server.Start()
}
