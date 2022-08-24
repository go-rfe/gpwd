package secrets

import (
	"errors"

	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Delete(id string) error {
	deleteSecretRequest := pb.DeleteSecretRequest{
		Secret: &pb.Secret{
			ID: id,
			Status: &pb.Status{
				Synced: false,
			},
		},
	}
	resp, err := c.grpc.DeleteSecret(c.ctx, &deleteSecretRequest)
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}

	return nil
}
