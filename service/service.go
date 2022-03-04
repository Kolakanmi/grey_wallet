package service

import (
	"context"
	"errors"
	"log"

	"github.com/Kolakanmi/grey_wallet/model"
	"github.com/Kolakanmi/grey_wallet/pkg/apperror"
	"github.com/Kolakanmi/grey_wallet/repository"
)

type IWalletService interface {
	Get(ctx context.Context) (*model.Wallet, error)
	Update(ctx context.Context, amount float64) (*model.Wallet, error)
}

type WalletService struct {
	walletRepository repository.IWalletRepository
}

func NewWalletService(walletRepository repository.IWalletRepository) *WalletService {
	return &WalletService{walletRepository: walletRepository}
}

func (t *WalletService) Get(ctx context.Context) (*model.Wallet, error) {
	wallet, err := t.walletRepository.Get(ctx)
	if err != nil {
		log.Printf("error: %v \n", err)
		return nil, err
	}
	return wallet, nil
}
func (t *WalletService) Update(ctx context.Context, amount float64) (*model.Wallet, error) {
	wallet, err := t.walletRepository.UpdateOrCreate(ctx, amount)
	if err != nil {
		log.Printf("error: %v \n", err)
		if errors.Is(err, repository.ErrInsufficientBalance) {
			return nil, apperror.BadRequestError("Insufficient balance")
		}
		return nil, apperror.CouldNotCompleteRequest()
	}
	return wallet, nil
}
