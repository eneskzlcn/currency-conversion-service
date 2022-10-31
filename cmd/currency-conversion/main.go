package main

import (
	"fmt"
	"github.com/eneskzlcn/currency-conversion-service/app/auth"
	"github.com/eneskzlcn/currency-conversion-service/app/conversion"
	"github.com/eneskzlcn/currency-conversion-service/app/exchange"
	"github.com/eneskzlcn/currency-conversion-service/app/wallet"
	"github.com/eneskzlcn/currency-conversion-service/config"
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

	authRepository := auth.NewRepository(db, logger)
	authService := auth.NewService(configs.Jwt, authRepository, logger)
	authHandler := auth.NewHandler(authService, logger)
	authGuard := auth.NewGuard(authService, logger)

	walletRepository := wallet.NewRepository(db, logger)
	walletService := wallet.NewService(walletRepository, logger)
	walletHandler := wallet.NewHandler(walletService, authGuard, logger)

	exchangeRepository := exchange.NewRepository(db, logger)
	exchangeService := exchange.NewService(exchangeRepository, logger)
	exchangeHandler := exchange.NewHandler(exchangeService, authGuard, logger)

	conversionService := conversion.NewService(walletService, logger)
	conversionHandler := conversion.NewHandler(conversionService, authGuard, logger)

	server := server.New([]server.Handler{
		authHandler,
		walletHandler,
		exchangeHandler,
		conversionHandler,
	}, configs.Server, logger)
	logger.Sync()
	return server.Start()
}
