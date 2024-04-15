package config

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	HttpServer *ConfigHttpServer
	GrpcServer *ConfigGrpcServer
	Logger     *ConfigLogger
	DataBase   *ConfigDB
}

type ConfigSheduler struct {
	GrpcServer *ConfigGrpcServer
	Broker     *ConfigBroker
	Logger     *ConfigLogger
	Exchange   *ExchangeBroker
	DataBase   *ConfigDB
}

type ExchangeBroker struct {
	Type       string
	Name       string
	Reliable   bool
	RoutingKey string
}

type ConfigHttpServer struct {
	Host    string
	Port    string
	Timeout time.Duration
}

type ConfigGrpcServer struct {
	Host    string
	Port    string
	Timeout time.Duration
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

type ConfigBroker struct {
	Host     string
	Port     string
	User     string
	Password string
}

func NewConfigApp(pahToFile string) Config {
	viper.AddConfigPath(pahToFile)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	configLog := ConfigLogger{}
	configDB := ConfigDB{
		NameDB:   viper.Get("nameDB").(string),
		Host:     viper.Get("hostDB").(string),
		Port:     viper.Get("portDB").(string),
		User:     viper.Get("usesrDB").(string),
		Password: viper.Get("passwordDB").(string),
	}

	configGrpcServer := ConfigGrpcServer{
		Host: viper.Get("hostServerGrpc").(string),
		Port: viper.Get("portServerGrpc").(string),
	}
	configHttpServer := ConfigHttpServer{
		Host: viper.Get("hostServerHttp").(string),
		Port: viper.Get("portServerHttp").(string),
	}

	level, err := zapcore.ParseLevel(viper.Get("Level").(string))
	if err != nil {
		configLog = ConfigLogger{zapcore.DebugLevel, viper.Get("logEncoding").(string)}
	} else {
		configLog = ConfigLogger{level, viper.Get("logEncoding").(string)}
	}

	return Config{
		HttpServer: &configHttpServer,
		GrpcServer: &configGrpcServer,
		Logger:     &configLog,
		DataBase:   &configDB,
	}
}

func NewConfigSheduler(pahToFile string) ConfigSheduler {
	viper.AddConfigPath(pahToFile)
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	configLog := ConfigLogger{}
	configDB := ConfigDB{
		NameDB:   viper.Get("nameDB").(string),
		Host:     viper.Get("hostDB").(string),
		Port:     viper.Get("portDB").(string),
		User:     viper.Get("usesrDB").(string),
		Password: viper.Get("passwordDB").(string),
	}
	configBroker := ConfigBroker{
		Host:     viper.Get("hostBroker").(string),
		Port:     viper.Get("portBroker").(string),
		User:     viper.Get("usesrBroker").(string),
		Password: viper.Get("passwordBroker").(string),
	}
	exchange := ExchangeBroker{
		Type:       viper.Get("exchangeType").(string),
		Name:       viper.Get("exchangeName").(string),
		Reliable:   viper.Get("exchangeReliable").(bool),
		RoutingKey: viper.Get("exchangeRoutingKey").(string),
	}
	configGrpcServer := ConfigGrpcServer{
		Host: viper.Get("hostServerGrpc").(string),
		Port: viper.Get("portServerGrpc").(string),
	}

	level, err := zapcore.ParseLevel(viper.Get("Level").(string))
	if err != nil {
		configLog = ConfigLogger{zapcore.DebugLevel, viper.Get("logEncoding").(string)}
	} else {
		configLog = ConfigLogger{level, viper.Get("logEncoding").(string)}
	}

	return ConfigSheduler{
		GrpcServer: &configGrpcServer,
		Exchange:   &exchange,
		Broker:     &configBroker,
		Logger:     &configLog,
		DataBase:   &configDB,
	}
}
