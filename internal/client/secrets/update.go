package secrets

import (
	"errors"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
)

func (c *client) Update(id string, data []byte, labels []string) (string, error) {
	labelsMap, err := constructLabels(labels)
	if err != nil {
		return "", err
	}

	secret := pb.Secret{
		ID:        id,
		Data:      data,
		Labels:    labelsMap,
		UpdatedAt: timestamppb.Now(),
		Status: &pb.Status{
			Synced: false,
		},
	}

	updateSecretRequest := pb.UpdateSecretRequest{
		Secret: &secret,
	}
	resp, err := c.grpc.UpdateSecret(c.ctx, &updateSecretRequest)
	if err != nil {
		return "", err
	}
	if resp.GetError() != "" {
		return "", errors.New(resp.GetError())
	}

	return resp.GetId(), nil
}
