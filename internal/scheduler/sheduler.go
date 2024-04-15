package scheduler

import (
	"context"
	"errors"

	"github.com/SashaMelva/calendar_service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func NewPublisher(body []byte,
	ctx context.Context, connection *amqp.Connection, confExchange *config.ExchangeBroker, log *zap.SugaredLogger, conf *config.ConfigBroker) error {

	log.Info("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return errors.New("Channel:  " + err.Error())
	}

	log.Info("got Channel, declaring %q Exchange (%q)", confExchange.Type, confExchange.Name)
	if err := channel.ExchangeDeclare(
		confExchange.Name, // name
		confExchange.Type, // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // noWait
		nil,               // arguments
	); err != nil {
		return errors.New("Exchange Declare: " + err.Error())
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if confExchange.Reliable {
		log.Info("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return errors.New("Channel could not be put into confirm mode: " + err.Error())
		}

		confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

		defer confirmOne(confirms, log)
	}

	log.Info("declared Exchange, publishing %dB body (%q)", len(body), body)
	if err = channel.PublishWithContext(
		ctx,
		confExchange.Name,       // publish to an exchange
		confExchange.RoutingKey, // routing to 0 or more queues
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return errors.New("Exchange Publish:  " + err.Error())
	}

	return nil
}

func confirmOne(confirms <-chan amqp.Confirmation, log *zap.SugaredLogger) {
	log.Info("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		log.Info("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Info("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
