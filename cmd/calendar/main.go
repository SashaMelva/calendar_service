package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/SashaMelva/calendar_service/internal/app"
	config "github.com/SashaMelva/calendar_service/internal/config"
	logger "github.com/SashaMelva/calendar_service/internal/logger"
	internalhttp "github.com/SashaMelva/calendar_service/internal/server/http"
	memorystorage "github.com/SashaMelva/calendar_service/internal/storage/memory"
	sqlstorage "github.com/SashaMelva/calendar_service/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../", "Path to configuration file")
}

func main() {
	fmt.Println("!231321")
	//flag.Parse()

	// if flag.Arg(0) == "version" {
	// 	printVersion()
	// 	return
	// }

	config := config.NewConfig(configFile)
	log := logger.NewLogger(config.Logger)
	//Соединение с бд
	connection := sqlstorage.New(config.DataBase, log)
	//Событие
	memstorage := memorystorage.New(connection.StorageDb)
	calendar := app.New(log, memstorage, config.App)

	server := internalhttp.NewServer(log, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
	}()

	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
