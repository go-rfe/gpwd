package accounts

import (
	"errors"

	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Update(server, username string, password []byte) error {
	account := pb.Account{
		ServerAddress: server,
		UserName:      username,
		UserPassword:  password,
	}

	updateSecretRequest := pb.UpdateAccountRequest{
		Account: &account,
	}
	resp, err := c.grpc.UpdateAccount(c.ctx, &updateSecretRequest)
	if err != nil {
		return err
	}
	if resp.GetError() != "" {
		return errors.New(resp.GetError())
	}

	return nil
}
