package repository

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/Kolakanmi/grey_wallet/model"
	"github.com/Kolakanmi/grey_wallet/pkg/utils"
	"github.com/Kolakanmi/grey_wallet/pkg/uuid"
)

type IWalletRepository interface {
	Create(ctx context.Context, amount float64) (*model.Wallet, error)
	Get(ctx context.Context) (*model.Wallet, error)
	UpdateOrCreate(ctx context.Context, amount float64) (*model.Wallet, error)
}

type WalletRepository struct {
	db *sql.DB
}

var ErrInsufficientBalance = errors.New("insufficient balance")

var (
	getStatement = `
		SELECT id, created_at, updated_at, balance FROM kola_wallets 
		where deleted_at is null
	`
	updateStatement = `
		UPDATE kola_wallets SET balance = balance + $1, updated_at = $2 RETURNING balance, updated_at;
	`
	createStatement = `
		INSERT INTO kola_wallets (id, created_at, updated_at, balance)
		VALUES ($1, $2, $3, $4)
	`
)

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (w *WalletRepository) Create(ctx context.Context, amount float64) (*model.Wallet, error) {
	timeNow := utils.TimeNow()
	id := uuid.New()

	_, err := w.db.ExecContext(ctx, createStatement, id, timeNow, timeNow, amount)
	if err != nil {
		return nil, err
	}
	return &model.Wallet{Base: model.Base{ID: id, CreatedAt: timeNow, UpdatedAt: timeNow}, Balance: amount}, nil
}

func (w *WalletRepository) Get(ctx context.Context) (*model.Wallet, error) {
	row := w.db.QueryRowContext(ctx, getStatement)
	var wallet model.Wallet
	err := row.Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return w.Create(ctx, 0)
		}
		return nil, err
	}
	return &wallet, nil
}

func (w *WalletRepository) UpdateOrCreate(ctx context.Context, amount float64) (*model.Wallet, error) {
	row := w.db.QueryRowContext(ctx, getStatement)
	var wallet model.Wallet
	err := row.Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt, &wallet.Balance)
	if err == nil {
		if amount < 0 && wallet.Balance < math.Abs(amount) {
			return nil, ErrInsufficientBalance
		}
		//update balance
		err = w.db.QueryRowContext(ctx, updateStatement, amount, time.Now()).
			Scan(&wallet.Balance, &wallet.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &wallet, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	if amount < 0 {
		return nil, ErrInsufficientBalance
	}
	//if no row, create wallet record
	return w.Create(ctx, amount)
}
