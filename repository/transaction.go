package repository

import (
	"context"
	"database/sql"

	"github.com/Kolakanmi/grey_wallet/model"
	"github.com/Kolakanmi/grey_wallet/pkg/utils"
	"github.com/Kolakanmi/grey_wallet/pkg/uuid"
)

type ITransactionRepository interface {
	Create(ctx context.Context, txn *model.Transaction) (string, error)
	GetAll(ctx context.Context) ([]model.Transaction, error)
	UpdateStatus(ctx context.Context, txnID string, status string) error
}

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (t *TransactionRepository) Create(ctx context.Context, txn *model.Transaction) (string, error) {
	timeNow := utils.TimeNow()
	id := uuid.New()
	statement := `
		INSERT INTO kola_transactions (id, created_at, updated_at, amount, status)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := t.db.ExecContext(ctx, statement, id, timeNow, timeNow, txn.Amount, model.Pending)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *TransactionRepository) GetAll(ctx context.Context) ([]model.Transaction, error) {
	query := `
		SELECT id, created_at, updated_at, amount, status FROM kola_transactions 
		where status = $1 and deleted_at is null
	`
	rows, err := t.db.QueryContext(ctx, query, model.Success)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []model.Transaction{}
	for rows.Next() {
		var txn model.Transaction
		err := rows.Scan(&txn.ID, &txn.CreatedAt, &txn.UpdatedAt, &txn.Amount, &txn.Status)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}

func (t *TransactionRepository) UpdateStatus(ctx context.Context, txnID string, status string) error {
	timeNow := utils.TimeNow()
	statement := `
		UPDATE kola_transactions SET status = $1, updated_at = $2 WHERE id = $3;
	`
	_, err := t.db.ExecContext(ctx, statement, status, timeNow, txnID)
	if err != nil {
		return err
	}
	return nil
}
