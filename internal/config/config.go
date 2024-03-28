package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	App      *ConfigApp
	Logger   *ConfigLogger
	DataBase *ConfigDB
}

type ConfigApp struct {
	Host string
	Port string
}

type ConfigLogger struct {
	Level       zapcore.Level
	LogEncoding string `required:"true"`
}

type ConfigDB struct {
	NameDB   string
	Host     string
	Port     string
	User     string
	Password string
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
	confDB := ConfigDB{
		NameDB:   viper.Get("nameDB").(string),
		Host:     viper.Get("hostDB").(string),
		Port:     viper.Get("portDB").(string),
		User:     viper.Get("usesrDB").(string),
		Password: viper.Get("passwordDB").(string),
	}
	confApp := ConfigApp{
		Host: viper.Get("hostServer").(string),
		Port: viper.Get("portServer").(string),
	}

	level, err := zapcore.ParseLevel(viper.Get("Level").(string))
	if err != nil {
		confLog = ConfigLogger{zapcore.DebugLevel, viper.Get("logEncoding").(string)}
	} else {
		confLog = ConfigLogger{level, viper.Get("logEncoding").(string)}
	}

	return Config{
		App:      &confApp,
		Logger:   &confLog,
		DataBase: &confDB,
	}
}
