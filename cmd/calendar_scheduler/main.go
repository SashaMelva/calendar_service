package main

import (
	"context"
	"encoding/json"
	"flag"
	"time"

	"github.com/SashaMelva/calendar_service/internal/config"
	"github.com/SashaMelva/calendar_service/internal/logger"
	"github.com/SashaMelva/calendar_service/internal/scheduler"
	internalgrpc "github.com/SashaMelva/calendar_service/internal/server/grpc"
	proto "github.com/SashaMelva/calendar_service/internal/server/grpc/gen"
	amqp "github.com/rabbitmq/amqp091-go"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../", "Path to configuration file")
}

func main() {
	//flag.Parse()

	// if flag.Arg(0) == "version" {
	// 	printVersion()
	// 	return
	// }

	config := config.NewConfigSheduler(configFile)
	log := logger.NewLogger(config.Logger)

	connection, err := amqp.Dial("amqp://" + config.Broker.User + ":" + config.Broker.Password + "@" + config.Broker.Host + ":" + config.Broker.Port + "/") // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	defer func() {
		_ = connection.Close() // Закрываем подключение в случае удачной попытки
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	client := internalgrpc.NewGRPCClient(config.GrpcServer, log)
	event, err := client.GetEventById(ctx, &proto.EventId{Id: 1})

	if err != nil {
		log.Error(err)
	}

	jsonBytes, err := json.Marshal(event)
	if err != nil {
		log.Error("Ошибка при сериализации в JSON:", err)
	}

	err = scheduler.NewPublisher(jsonBytes, ctx, connection, config.Exchange, log, config.Broker)

	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	log.Info("sheduler is running...")

}
