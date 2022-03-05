package adapter

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "github.com/Kolakanmi/grey_wallet/pkg/grpc/transaction"
)

func NewClientConnection(config *Config) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", config.Address, config.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewClient(conn *grpc.ClientConn) proto.WalletClient {
	wc := proto.NewWalletClient(conn)
	return wc
}

func CloseConnection(conn *grpc.ClientConn) {
	conn.Close()
}
