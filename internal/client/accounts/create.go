package accounts

import (
	"errors"

	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Create(server string, username string, password []byte) (string, error) {
	account := pb.Account{
		ServerAddress: server,
		UserName:      username,
		UserPassword:  password,
	}

	createAccountRequest := pb.CreateAccountRequest{
		Account: &account,
	}
	resp, err := c.grpc.CreateAccount(c.ctx, &createAccountRequest)
	if err != nil {
		return "", err
	}
	if resp.GetError() != "" {
		return "", errors.New(resp.GetError())
	}

	return resp.GetId(), nil
}
