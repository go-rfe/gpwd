package local

import (
	"context"
	"errors"

	pb "github.com/go-rfe/gpwd/internal/proto"
)

var ErrNoSecretFound = errors.New("no secrets found with provided id")

type Secrets interface {
	CreateSecret(ctx context.Context, secret *pb.Secret) (string, error)
	ListSecrets(ctx context.Context) ([]*pb.Secret, error)
	GetSecret(ctx context.Context, id string) (*pb.Secret, error)
	UpdateSecret(ctx context.Context, secret *pb.Secret) error
	DeleteSecret(ctx context.Context, secret *pb.Secret) error
}
