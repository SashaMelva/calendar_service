package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/SashaMelva/calendar_service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type MessageEvent struct {
	Id        uint32
	Title     string
	DateStart time.Time
}

func NewMessage(messge *MessageEvent) ([]byte, error) {
	if messge.Id == 0 || messge.Title == "" {
		return []byte(""), errors.New("Error Valid data for message event")
	}

	jsonBytes, err := json.Marshal(messge)
	if err != nil {
		return []byte(""), errors.New("Ошибка при сериализации в JSON:" + err.Error())
	}

	return jsonBytes, nil
}

func OpenConnection(log *zap.SugaredLogger, conf *config.ConfigBroker) *amqp.Connection {
	conn, err := amqp.Dial("amqp://" + conf.User + ":" + conf.Password + "@" + conf.Host + ":" + conf.Port + "/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatal("unable to open connect to RabbitMQ server. Error:" + err.Error())
	}

	return conn
}

func QueueDeclare(log *zap.SugaredLogger, channel *amqp.Channel, name string) {
	_, err := channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}
}

func PublishMessage(ctx context.Context, queueName string, body []byte, log *zap.SugaredLogger, channel *amqp.Channel) {
	if err := channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
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
		log.Fatal("failed to publish a message. Error:  " + err.Error())
	}
}
