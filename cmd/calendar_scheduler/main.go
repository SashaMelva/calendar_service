package main

import (
	"context"
	"flag"
	"time"

	"github.com/SashaMelva/calendar_service/internal/config"
	"github.com/SashaMelva/calendar_service/internal/logger"
	"github.com/SashaMelva/calendar_service/internal/rabbit"
	internalgrpc "github.com/SashaMelva/calendar_service/internal/server/grpc"
	proto "github.com/SashaMelva/calendar_service/internal/server/grpc/gen"
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

	conn := rabbit.OpenConnection(log, config.Broker)

	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки
	}()

	log.Info("got Connection, getting Channel")
	channel, err := conn.Channel()

	if err != nil {
		log.Fatal("Channel:  " + err.Error())
	}

	defer func() {
		_ = channel.Close() // Закрываем подключение в случае удачной попытки подключения
	}()

	rabbit.QueueDeclare(log, channel, "sendingEvents")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	client := internalgrpc.NewGRPCClient(config.GrpcServer, log)
	event, err := client.GetEventById(ctx, &proto.EventId{Id: 1})

	if err != nil {
		log.Error(err)
	}

	message, err := rabbit.NewMessage(&rabbit.MessageEvent{
		Id:        event.Event.Id,
		Title:     event.Event.Title,
		DateStart: event.Event.DateTimeEnd.AsTime(),
	})

	if err != nil {
		log.Fatal(err)
	}

	rabbit.PublishMessage(ctx, "sendingEvents", message, log, channel)

	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	log.Info("sheduler is running...")

}
