package secrets

import (
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Create(data []byte, labels []string) (string, error) {
	labelsMap, err := constructLabels(labels)
	if err != nil {
		return "", err
	}

	secret := pb.Secret{
		Data:      data,
		Labels:    labelsMap,
		CreatedAt: timestamppb.Now(),
		Status: &pb.Status{
			Synced: false,
		},
	}

	createSecretRequest := pb.CreateSecretRequest{
		Secret: &secret,
	}
	resp, err := c.grpc.CreateSecret(c.ctx, &createSecretRequest)
	if err != nil {
		return "", err
	}
	if resp.GetError() != "" {
		return "", errors.New(resp.GetError())
	}

	return resp.GetId(), nil
}
