package cloud

import (
	"context"
	"errors"

	pb "github.com/go-rfe/gpwd/internal/proto"
)

var (
	ErrAccountExists   = errors.New("user account already exists")
	ErrAccountNotFound = errors.New("user account not found")
)

type Accounts interface {
	CreateAccount(ctx context.Context, auth *pb.Auth) error
	GetByName(ctx context.Context, username string) (*pb.Auth, error)
}
