package accounts

import (
	"errors"

	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Delete() error {
	deleteAccountRequest := pb.DeleteAccountRequest{}
	resp, err := c.grpc.DeleteAccount(c.ctx, &deleteAccountRequest)
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}

	return nil
}
