package wallet

import (
	"context"
	"encoding/json"
	"github.com/eneskzlcn/currency-conversion-service/app/message"
	"github.com/eneskzlcn/currency-conversion-service/config"
	"go.uber.org/zap"
)

const (
	RabbitmqConsumerName = "wallet-service-currency-converted-consumer"
)

type RabbitmqClient interface {
	Consume(messageReceived chan []byte, consumer string, queue string)
}

type rabbitmqConsumer struct {
	client  RabbitmqClient
	logger  *zap.SugaredLogger
	service Service
	config  config.RabbitMQ
}

func NewRabbitmqConsumer(client RabbitmqClient, logger *zap.SugaredLogger,
	service Service, config config.RabbitMQ) *rabbitmqConsumer {
	return &rabbitmqConsumer{client: client, logger: logger,
		service: service, config: config}
}

func (r *rabbitmqConsumer) ConsumeCurrencyConvertedQueue() {
	onMessageReceived := make(chan []byte, 0)
	go r.client.Consume(onMessageReceived, RabbitmqConsumerName, r.config.CurrencyConvertedQueue)
	var forever chan struct{}
	for messageBytes := range onMessageReceived {
		r.logger.Debugf("Consumed a message from '%s' queue", r.config.CurrencyConvertedQueue)
		var currencyConvertedMessage message.CurrencyConvertedMessage
		if err := json.Unmarshal(messageBytes, &currencyConvertedMessage); err != nil {
			r.logger.Error(err)
			continue
		}
		err := r.service.TransferBalancesBetweenUserWallets(context.Background(), currencyConvertedMessage)
		if err != nil {
			r.logger.Error(err)
		}
	}
	<-forever
}
