package main

import (
	"flag"

	"github.com/SashaMelva/calendar_service/internal/app"
	"github.com/SashaMelva/calendar_service/internal/config"
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
	//Соединение с бд
	connection := sqlstorage.New(config.DataBase, log)
	//Событие
	memstorage := memorystorage.New(connection.StorageDb)
	calendar := app.New(log, memstorage)

	grpcServer := internalgrpc.NewGRPCServer(log, calendar)
	internalgrpc.ListenServer(grpcServer, config.GrpcServer, log)

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
