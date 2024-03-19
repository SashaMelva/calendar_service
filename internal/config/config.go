package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Url    string
	Logger *ConfigLogger
}

type ConfigLogger struct {
	Level       zapcore.Level
	LogEncoding string `required:"true"`
}

func NewConfig(pahToFile string) Config {
	viper.AddConfigPath(pahToFile)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	confLog := ConfigLogger{}

	level, err := zapcore.ParseLevel(viper.Get("Level").(string))
	if err != nil {
		confLog = ConfigLogger{zapcore.DebugLevel, viper.Get("logEncoding").(string)}
	} else {
		confLog = ConfigLogger{level, viper.Get("logEncoding").(string)}
	}

	conf := Config{viper.Get("Url").(string), &confLog}

	return conf
}
