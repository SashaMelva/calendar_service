package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/SashaMelva/calendar_service/internal/app"
	config "github.com/SashaMelva/calendar_service/internal/config"
	logger "github.com/SashaMelva/calendar_service/internal/logger"
	internalhttp "github.com/SashaMelva/calendar_service/internal/server/http"
	memorystorage "github.com/SashaMelva/calendar_service/internal/storage/memory"
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

	config := config.NewConfig(configFile)
	logg := logger.NewLogger(config.Logger)

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
