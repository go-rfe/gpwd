package secrets

import (
	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) List() ([]*pb.Secret, error) {
	resp, err := c.grpc.ListSecrets(c.ctx, &pb.ListSecretsRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetSecrets(), nil
}
