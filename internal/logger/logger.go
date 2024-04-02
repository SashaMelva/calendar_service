package logger

import (
	"fmt"
	"reflect"
	"time"

	config "github.com/SashaMelva/calendar_service/internal/config"
	"go.uber.org/zap"
)

func NewLogger(conf *config.ConfigLogger) *zap.SugaredLogger {
	fileName := time.Now()
	fmt.Println(conf.Level, reflect.TypeOf(conf.Level))

	logConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(conf.Level),
		DisableCaller:    true,
		Development:      true,
		Encoding:         conf.LogEncoding,
		OutputPaths:      []string{"stdout", "../../filelog/" + fileName.Format("01-02-2006") + ".log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}

	logger := zap.Must(logConfig.Build()).Sugar()

	logger.Info("Started")
	logger.Debug("Debug mode enabled")
	return logger
}
