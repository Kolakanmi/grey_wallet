package adapter

import (
	"github.com/Kolakanmi/grey_wallet/pkg/envconfig"
)

type Config struct {
	Address string `envconfig:"GRPC_ADDRESS" default:"localhost"`
	Port    int    `envconfig:"GRPC_PORT" default:"9090"`
}

func LoadConfig() *Config {
	var conf Config
	envconfig.Load(&conf)
	return &conf
}
