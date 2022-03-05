package service

import (
	"context"
	"errors"
	"log"

	"github.com/Kolakanmi/grey_wallet/model"
	"github.com/Kolakanmi/grey_wallet/pkg/apperror"
	"github.com/Kolakanmi/grey_wallet/repository"

	proto "github.com/Kolakanmi/grey_wallet/pkg/grpc/transaction"
)

type IWalletService interface {
	Get(ctx context.Context) (*model.Wallet, error)
	Update(ctx context.Context, amount float64) (*model.Wallet, error)
}

type WalletService struct {
	proto.UnimplementedWalletServer
	walletRepository repository.IWalletRepository
}

func NewWalletService(walletRepository repository.IWalletRepository) *WalletService {
	return &WalletService{walletRepository: walletRepository}
}

func (w *WalletService) Get(ctx context.Context) (*model.Wallet, error) {
	wallet, err := w.walletRepository.Get(ctx)
	if err != nil {
		log.Printf("error: %v \n", err)
		return nil, err
	}
	return wallet, nil
}
func (w *WalletService) Update(ctx context.Context, amount float64) (*model.Wallet, error) {
	wallet, err := w.walletRepository.UpdateOrCreate(ctx, amount)
	if err != nil {
		log.Printf("error: %v \n", err)
		if errors.Is(err, model.ErrInsufficientBalance) {
			return nil, apperror.BadRequestError("Insufficient balance")
		}
		return nil, apperror.CouldNotCompleteRequest()
	}
	return wallet, nil
}

func (w *WalletService) GetBalance(ctx context.Context, req *proto.GetBalanceRequest) (*proto.GetBalanceResponse, error) {
	wallet, err := w.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.GetBalanceResponse{Balance: wallet.Balance}, nil
}

func (w *WalletService) UpdateBalance(ctx context.Context, req *proto.UpdateBalanceRequest) (*proto.UpdateBalanceResponse, error) {
	wallet, err := w.Update(ctx, req.Amount)
	if err != nil {
		return nil, err
	}
	return &proto.UpdateBalanceResponse{Balance: wallet.Balance}, nil
}
