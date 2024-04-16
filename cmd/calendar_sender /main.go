package main

import (
	"flag"
	"time"

	"github.com/SashaMelva/calendar_service/internal/config"
	"github.com/SashaMelva/calendar_service/internal/consumer"
	"github.com/SashaMelva/calendar_service/internal/logger"
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

	c, err := consumer.NewConsumer(&config.Exchange, &config.Broker)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if *lifetime > 0 {
		log.Printf("running for %s", *lifetime)
		time.Sleep(*lifetime)
	} else {
		log.Printf("running forever")
		select {}
	}

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}

	// httpServer := internalhttp.NewServer(log, calendar, config.HttpServer)

	// ctx, cancel := signal.NotifyContext(context.Background(),
	// 	syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// defer cancel()

	// go func() {
	// 	<-ctx.Done()

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	// 	defer cancel()

	// 	if err := httpServer.Stop(ctx); err != nil {
	// 		log.Error("failed to stop http server: " + err.Error())
	// 	}
	// }()

	log.Info("calendar is running...")

	// if err := httpServer.Start(ctx); err != nil {
	// 	log.Error("failed to start http server: " + err.Error())
	// 	cancel()
	// 	os.Exit(1) //nolint:gocritic
	// }
}
