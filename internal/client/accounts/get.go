package accounts

import (
	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Get() (*pb.Account, error) {
	resp, err := c.grpc.GetAccount(c.ctx, &pb.GetAccountRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetAccount(), nil
}
