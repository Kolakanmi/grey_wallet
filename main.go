package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/Kolakanmi/grey_wallet/adapter"
	"github.com/Kolakanmi/grey_wallet/pkg/database"
	"github.com/Kolakanmi/grey_wallet/pkg/envconfig"
	"github.com/Kolakanmi/grey_wallet/repository"
	"github.com/Kolakanmi/grey_wallet/service"
	"google.golang.org/grpc"

	proto "github.com/Kolakanmi/grey_wallet/pkg/grpc/transaction"
)

func main() {
	err := envconfig.SetEnvFromConfig(".env")
	if err != nil {
		log.Println("env config load err: ", err)
	}

	grpcServer := grpc.NewServer()

	dbConf := database.LoadConfig()
	db, err := database.ConnectDB(dbConf)
	if err != nil {
		log.Printf("db error %v", err)
		return
	}

	repo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(repo)

	// walletServer := adapter.NewWalletServer(walletService)

	adapterConfig := adapter.LoadConfig()

	proto.RegisterWalletServer(grpcServer, walletService)

	address := fmt.Sprintf("%s:%d", adapterConfig.Address, adapterConfig.Port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("listen error: ", err)
		return
	}

	log.Printf("Wallet Server is listening on: %s", address)
	//Run server in goroutine to implement graceful shutdown.
	go func() {
		grpcServer.Serve(l)
		if err != nil {
			log.Printf("listen: %v\n", err)
		}
	}()

	//Graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	//Block until signal is received
	<-signals
	log.Print("Wallet server shutting down!!!")
	grpcServer.GracefulStop()
}
