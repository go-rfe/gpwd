package secrets

import (
	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Get(id string) (*pb.Secret, error) {
	resp, err := c.grpc.GetSecret(c.ctx, &pb.GetSecretRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetSecret(), nil
}
