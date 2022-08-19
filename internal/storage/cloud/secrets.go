package cloud

import (
	"context"

	pb "github.com/go-rfe/gpwd/internal/proto"
)

type Secrets interface {
	CreateSecrets(ctx context.Context, auth *pb.Auth, secrets []*pb.Secret) error
	DeleteSecrets(ctx context.Context, auth *pb.Auth, secrets []*pb.Secret) error
	UpdateSecrets(ctx context.Context, auth *pb.Auth, secrets []*pb.Secret) error
	ListSecrets(ctx context.Context, auth *pb.Auth) ([]*pb.Secret, error)
}
