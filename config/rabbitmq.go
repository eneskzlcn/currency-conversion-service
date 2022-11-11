package config

type RabbitMQ struct {
	Username               string
	Password               string
	Host                   string
	Port                   string
	CurrencyConvertedQueue string
}

func (r RabbitMQ) Queues() []string {
	queues := make([]string, 0)
	queues = append(queues, r.CurrencyConvertedQueue)
	return queues
}
