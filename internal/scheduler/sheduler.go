package scheduler

import (
	"context"

	"github.com/SashaMelva/calendar_service/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func NewPublisher(body []byte,
	ctx context.Context, connection *amqp.Connection, confExchange *config.ExchangeBroker, log *zap.SugaredLogger, conf *config.ConfigBroker) error {

	log.Info("got Channel, declaring %q Exchange (%q)", confExchange.Type, confExchange.Name)

	return nil
}
