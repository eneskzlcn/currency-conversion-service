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
	"github.com/eneskzlcn/currency-conversion-service/rabbitmq"
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
	rabbitmqClient := rabbitmq.NewClient(appConfig.Rabbitmq, logger)

	authRepository := auth.NewPostgresRepository(db, logger)
	authService := auth.NewService(appConfig.Jwt, authRepository, logger)
	authHttpHandler := auth.NewHttpHandler(authService, logger)
	authGuard := auth.NewGuard(authService, logger)

	walletRepository := wallet.NewPostgresRepository(db, logger)
	walletService := wallet.NewService(walletRepository, logger)
	walletHttpHandler := wallet.NewHttpHandler(walletService, authGuard, logger)
	walletRabbitmqConsumer := wallet.NewRabbitmqConsumer(&wallet.ConsumerOptions{
		Client:  rabbitmqClient,
		Logger:  logger,
		Service: walletService,
		Config:  appConfig.Rabbitmq,
	})

	exchangeRepository := exchange.NewPostgresRepository(db, logger)
	exchangeService := exchange.NewService(exchangeRepository, logger)
	exchangeHttpHandler := exchange.NewHttpHandler(exchangeService, authGuard, logger)

	conversionRepository := conversion.NewPostgresRepository(db, logger)
	conversionRabbitmqProducer := conversion.NewRabbitMqProducer(rabbitmqClient, appConfig.Rabbitmq)

	conversionService := conversion.NewService(&conversion.ServiceOptions{
		WalletService:    walletService,
		Logger:           logger,
		Repository:       conversionRepository,
		RabbitmqProducer: conversionRabbitmqProducer,
	})
	conversionHttpHandler := conversion.NewHttpHandler(conversionService, authGuard, logger)

	server := server.New([]server.Handler{
		authHttpHandler,
		walletHttpHandler,
		exchangeHttpHandler,
		conversionHttpHandler,
	}, appConfig.Server, logger)

	go walletRabbitmqConsumer.ConsumeCurrencyConvertedQueue()

	print("im trying to something new!")
	return server.Start()
}
