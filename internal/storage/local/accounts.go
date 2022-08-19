package local

import (
	"context"

	pb "github.com/go-rfe/gpwd/internal/proto"
)

type Accounts interface {
	CreateAccount(ctx context.Context, secret *pb.Account) (string, error)
	GetAccount(ctx context.Context) (*pb.Account, error)
	UpdateAccount(ctx context.Context, secret *pb.Account) error
	DeleteAccount(ctx context.Context) error
}
