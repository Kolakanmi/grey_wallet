package service

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/Kolakanmi/grey_wallet/model"
	"github.com/Kolakanmi/grey_wallet/pkg/apperror"
	proto "github.com/Kolakanmi/grey_wallet/pkg/grpc/transaction"
	"github.com/Kolakanmi/grey_wallet/repository/mock"
	mockConn "github.com/Kolakanmi/grey_wallet/service/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getService() *WalletService {
	repo := mock.NewMockRepo()

	// client := mockConn.NewMockClient()
	return NewWalletService(repo)
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		res  *model.Wallet
		err  error
	}{
		{
			name: "Get balance",
			res:  &model.Wallet{Balance: 100},
			err:  nil,
		},
	}
	t.Run("Get", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := getService()
				got, err := s.Get(context.Background())
				if err != nil {
					if tt.err == nil {
						t.Errorf("Expected no error, but got %v", err)
					}
					if err.Error() != tt.err.Error() {
						t.Errorf("Expected error to be %v, but got %v", tt.err, err)
					}
				}
				if tt.res != nil && got.Balance != tt.res.Balance {
					t.Errorf("Expected balance to be %v, but got %v", tt.res.Balance, got.Balance)
				}
			})
		}
	})
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		res    *model.Wallet
		err    error
	}{
		{
			name:   "When amount is greater than 0",
			amount: 10.0,
			res:    &model.Wallet{Balance: 110},
			err:    nil,
		},
		{
			name:   "When amount is less than 0",
			amount: -10.0,
			res:    &model.Wallet{Balance: 90},
			err:    nil,
		},
		{
			name:   "When amount is less than balance",
			amount: -110.0,
			res:    nil,
			err:    apperror.BadRequestError("Insufficient balance"),
		},
	}
	t.Run("Update", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := getService()
				got, err := s.Update(context.Background(), tt.amount)
				if err != nil {
					if tt.err == nil {
						t.Errorf("Expected no error, but got %v", err)
					}
					if err.Error() != tt.err.Error() {
						t.Errorf("Expected error to be %v, but got %v", tt.err, err)
					}
				}
				if tt.res != nil && got.Balance != tt.res.Balance {
					t.Errorf("Expected balance to be %v, but got %v", tt.res.Balance, got.Balance)
				}
			})
		}
	})
}

func TestGetBalance(t *testing.T) {
	tests := []struct {
		name    string
		balance float64
		err     error
	}{
		{
			name:    "Get balance",
			balance: 100.0,
			err:     nil,
		},
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(mockConn.Dialer()))
	if err != nil {
		log.Fatal(err)
	}
	wc := proto.NewWalletClient(conn)
	t.Run("GetBalance", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// s := getService()
				req := &proto.GetBalanceRequest{}
				got, err := wc.GetBalance(ctx, req)
				if err != nil {
					if tt.err == nil {
						t.Errorf("Expected no error, but got %v", err)
					}
					if err.Error() != tt.err.Error() {
						t.Errorf("Expected error to be %v, but got %v", tt.err, err)
					}
				}
				if got.Balance != tt.balance {
					t.Errorf("Expected balance to be %v, but got %v", tt.balance, got.Balance)
				}
			})
		}
	})
}

func TestUpdateBalance(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		balance float64
		err     error
	}{
		{
			name:    "When amount is greater than 0",
			amount:  10.0,
			balance: 110.0,
			err:     nil,
		},
		{
			name:    "When amount is less than 0",
			amount:  -10.0,
			balance: 90.0,
			err:     nil,
		},
		{
			name:    "When amount is less than balance",
			amount:  -110.0,
			balance: 100.0,
			err:     apperror.BadRequestError("insufficient balance"),
		},
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(mockConn.Dialer()))
	if err != nil {
		log.Fatal(err)
	}
	wc := proto.NewWalletClient(conn)
	t.Run("UpdateBalance", func(t *testing.T) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// s := getService()
				req := &proto.UpdateBalanceRequest{Amount: tt.amount}
				got, err := wc.UpdateBalance(ctx, req)
				// fmt.Fprintf("got: %+v", got)
				if err != nil {
					if tt.err == nil {
						t.Errorf("Expected no error, but got %v", err)
					}
					if !strings.Contains(err.Error(), tt.err.Error()) {
						t.Errorf("Expected error to be %v, but got %v", tt.err, err)
					}
				}
				if err == nil && got.Balance != tt.balance {
					t.Errorf("Expected balance to be %v, but got %v", tt.balance, got.Balance)
				}
			})
		}
	})
}
