package service

import (
	"context"
	"log"

	"github.com/Kolakanmi/grey_wallet/model"
	"github.com/Kolakanmi/grey_wallet/pkg/apperror"
	proto "github.com/Kolakanmi/grey_wallet/pkg/grpc/transaction"
	"github.com/Kolakanmi/grey_wallet/repository"
)

type ITransactionService interface {
	Credit(ctx context.Context, amount float64) (*TxnResponse, error)
	Debit(ctx context.Context, amount float64) (*TxnResponse, error)
	Balance(ctx context.Context) (*TxnResponse, error)
	GetAll(ctx context.Context) ([]model.Transaction, error)
}

type TransactionService struct {
	transactionRepository repository.ITransactionRepository
	walletClient          proto.WalletClient
}

func NewTransactionService(transactionRepository repository.ITransactionRepository, wc proto.WalletClient) *TransactionService {
	return &TransactionService{transactionRepository: transactionRepository, walletClient: wc}
}

func (t *TransactionService) Credit(ctx context.Context, amount float64) (*TxnResponse, error) {
	if amount < 0 {
		return nil, apperror.BadRequestError("Amount cannot be negative.")
	}
	id, err := t.transactionRepository.Create(ctx, &model.Transaction{
		Amount: amount,
		Status: model.Pending,
	})
	if err != nil {
		log.Printf("error: %v \n", err)
		return nil, apperror.CouldNotCompleteRequest()
	}

	wallet, err := t.walletClient.UpdateBalance(ctx, &proto.UpdateBalanceRequest{Amount: amount})
	if err != nil {
		log.Printf("error: %v \n", err)
		t.transactionRepository.UpdateStatus(ctx, id, model.Failed)
		return nil, err
	}

	t.transactionRepository.UpdateStatus(ctx, id, model.Success)

	return &TxnResponse{Balance: wallet.Balance}, nil
}

func (t *TransactionService) Debit(ctx context.Context, amount float64) (*TxnResponse, error) {
	response, err := t.walletClient.GetBalance(ctx, &proto.GetBalanceRequest{})
	if err != nil {
		log.Printf("error: %v \n", err)
		return nil, apperror.CouldNotCompleteRequest()
	}
	if response.Balance < amount {
		return nil, apperror.BadRequestError("Insufficient balance")
	}
	if amount > 0 {
		amount = -amount
	}
	id, err := t.transactionRepository.Create(ctx, &model.Transaction{
		Amount: amount,
		Status: model.Pending,
	})
	if err != nil {
		log.Printf("error: %v \n", err)
		return nil, apperror.CouldNotCompleteRequest()
	}

	wallet, err := t.walletClient.UpdateBalance(ctx, &proto.UpdateBalanceRequest{Amount: amount})
	if err != nil {
		log.Printf("error: %v \n", err)
		t.transactionRepository.UpdateStatus(ctx, id, model.Failed)
		return nil, err
	}

	t.transactionRepository.UpdateStatus(ctx, id, model.Success)

	return &TxnResponse{Balance: wallet.Balance}, nil
}

func (t *TransactionService) Balance(ctx context.Context) (*TxnResponse, error) {
	wallet, err := t.walletClient.GetBalance(ctx, &proto.GetBalanceRequest{})
	if err != nil {
		log.Printf("error: %v \n", err)
		return nil, err
	}

	return &TxnResponse{Balance: wallet.Balance}, nil
}

func (t *TransactionService) GetAll(ctx context.Context) ([]model.Transaction, error) {
	txns, err := t.transactionRepository.GetAll(ctx)
	if err != nil {
		log.Printf("error: %v \n", err)
		return nil, err
	}

	return txns, nil
}
