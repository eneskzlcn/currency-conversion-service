package conversion

import "github.com/eneskzlcn/currency-conversion-service/config"

type RabbitmqClient interface {
	PushMessage(message any, queue string) error
}

type rabbitmqProducer struct {
	client RabbitmqClient
	config config.RabbitMQ
}

func NewRabbitMqProducer(client RabbitmqClient, config config.RabbitMQ) *rabbitmqProducer {
	return &rabbitmqProducer{client: client, config: config}
}
func (r *rabbitmqProducer) PushConversionCreatedMessage(message CurrencyConvertedMessage) error {
	return r.client.PushMessage(message, r.config.CurrencyConvertedQueue)
}
