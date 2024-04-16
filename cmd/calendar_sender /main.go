package main

import (
	"flag"

	"github.com/SashaMelva/calendar_service/internal/config"
	"github.com/SashaMelva/calendar_service/internal/logger"
	"github.com/SashaMelva/calendar_service/internal/rabbit"
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

	config := config.NewConfigSender(configFile)
	log := logger.NewLogger(config.Logger)

	conn := rabbit.OpenConnection(log, config.Broker)

	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки подключения
	}()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. Error: %s", err)
	}

	defer func() {
		_ = channel.Close() // Закрываем подключение в случае удачной попытки подключения
	}()

	rabbit.QueueDeclare(log, channel, "sendingEvents")

	messages, err := channel.Consume(
		"sendingEvents", // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)

	if err != nil {
		log.Fatalf("failed to register a consumer. Error: %s", err)
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			log.Info(message.Body)
		}
	}()

	log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
