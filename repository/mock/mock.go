package mock

import (
	"context"
	"math"
	"time"

	"github.com/Kolakanmi/grey_wallet/model"
	"github.com/Kolakanmi/grey_wallet/pkg/uuid"
)

type MockRepo struct{}

func NewMockRepo() *MockRepo {
	return &MockRepo{}
}

// var ErrInsufficientBalance = errors.New("insufficient balance")

var timeNow, _ = time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")

var balance = 100.0

func getWithAmount(ctx context.Context, amount float64) (*model.Wallet, error) {
	wallet := &model.Wallet{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		},
		Balance: amount,
	}
	return wallet, nil
}

func (m *MockRepo) Get(ctx context.Context) (*model.Wallet, error) {
	return getWithAmount(ctx, balance)
}

func (m *MockRepo) Create(ctx context.Context, amount float64) (*model.Wallet, error) {
	return getWithAmount(ctx, amount)
}

func (m *MockRepo) UpdateOrCreate(ctx context.Context, amount float64) (*model.Wallet, error) {
	if amount < 0 && balance < math.Abs(amount) {
		return nil, model.ErrInsufficientBalance
	}
	return getWithAmount(ctx, balance+amount)
}
