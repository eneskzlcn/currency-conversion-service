package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eneskzlcn/currency-conversion-service/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type client struct {
	connection *amqp.Connection
	logger     *zap.SugaredLogger
}

func NewClient(config config.RabbitMQ, logger *zap.SugaredLogger) *client {
	if logger == nil {
		fmt.Println("given *zap.SugaredLogger is nil")
		return nil
	}
	con, err := amqp.Dial(createConnectionUrl(config))
	if err != nil {
		logger.Error("error occurred when connecting to rabbitmq server")
		return nil
	}
	ch, err := con.Channel()
	defer ch.Close()
	for _, queue := range config.Queues() {
		_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			logger.Error("error when declaring new queue")
			return nil
		}
	}

	return &client{connection: con, logger: logger}
}

func (c *client) PushMessage(message any, queue string) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	ch, err := c.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	context, cancelFn := context.WithTimeout(context.Background(), time.Second*40)
	defer cancelFn()
	err = ch.PublishWithContext(context,
		"", queue, false, false,
		amqp.Publishing{
			Headers:     nil,
			ContentType: "text/plain",
			Body:        messageBytes,
		})
	if err != nil {
		c.logger.Error("error occurred when publishing the message ")
		return err
	}
	return nil
}

func (c *client) Consume(messageReceived chan []byte, consumer string, queue string) {
	ch, err := c.connection.Channel()
	defer ch.Close()
	if err != nil {
		return
	}
	msgs, err := ch.Consume(
		queue,
		consumer,
		true,
		false,
		false,
		false,
		nil,
	)
	var forever chan struct{}
	go func() {
		for d := range msgs {
			messageReceived <- d.Body
		}
	}()
	<-forever
}

func createConnectionUrl(config config.RabbitMQ) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)
}
