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
	appConfig := config.New()
	appConfig.Print()
	logger, err := logger.NewZapLoggerForEnv(appConfig.AppEnv, 0)
	if err != nil {
		return err
	}

	db, err := postgres.New(appConfig.Db)
	if err != nil {
		return err
	}

	authRepository := auth.NewRepository(db, logger)
	authService := auth.NewService(appConfig.Jwt, authRepository, logger)
	authHttpHandler := auth.NewHttpHandler(authService, logger)
	authGuard := auth.NewGuard(authService, logger)

	walletRepository := wallet.NewRepository(db, logger)
	walletService := wallet.NewService(walletRepository, logger)
	walletHttpHandler := wallet.NewHttpHandler(walletService, authGuard, logger)

	exchangeRepository := exchange.NewRepository(db, logger)
	exchangeService := exchange.NewService(exchangeRepository, logger)
	exchangeHttpHandler := exchange.NewHttpHandler(exchangeService, authGuard, logger)

	conversionRepository := conversion.NewRepository(db, logger)
	conversionService := conversion.NewService(walletService, logger, conversionRepository)
	conversionHttpHandler := conversion.NewHttpHandler(conversionService, authGuard, logger)

	server := server.New([]server.Handler{
		authHttpHandler,
		walletHttpHandler,
		exchangeHttpHandler,
		conversionHttpHandler,
	}, appConfig.Server, logger)

	return server.Start()
}
