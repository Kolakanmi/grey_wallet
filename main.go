package main

import (
	"log"

	"github.com/Kolakanmi/grey_wallet/api"
	"github.com/Kolakanmi/grey_wallet/pkg/envconfig"
	"github.com/Kolakanmi/grey_wallet/pkg/http/server"
)

func main() {
	err := envconfig.SetEnvFromConfig(".env")
	if err != nil {
		log.Println("env config load err: ", err)
	}

	r, err := api.NewRouter()
	if err != nil {
		log.Printf("error : %v", err)
	}
	serverConf := server.LoadConfigFromEnv()
	server.ListenAndServe(*serverConf, r)
}
