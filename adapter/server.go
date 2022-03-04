package adapter

import (
	"context"

	proto "github.com/Kolakanmi/grey_wallet/pkg/grpc/transaction"
	"github.com/Kolakanmi/grey_wallet/service"
)

type WalletServer struct {
	service service.IWalletService
}

func NewWalletServer(service service.IWalletService) *WalletServer {
	return &WalletServer{service: service}
}

func (w *WalletServer) GetBalance(ctx context.Context, req *proto.GetBalanceRequest) (*proto.GetBalanceResponse, error) {
	wallet, err := w.service.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.GetBalanceResponse{Balance: wallet.Balance}, nil
}
func (w *WalletServer) UpdateBalance(ctx context.Context, req *proto.UpdateBalanceRequest) (*proto.UpdateBalanceResponse, error) {
	wallet, err := w.service.Update(ctx, req.Amount)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateBalanceResponse{Balance: wallet.Balance}, nil
}
