package accounts

import (
	"context"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
)

type client struct {
	grpc pb.AccountsClient
	ctx  context.Context
}

func NewAccountsClient(ctx context.Context, socket string) (*client, error) {
	grpcClient, err := getGRPCClient(ctx, socket)
	if err != nil {
		return nil, err
	}

	return &client{grpc: grpcClient, ctx: ctx}, nil
}

func getGRPCClient(ctx context.Context, socket string) (pb.AccountsClient, error) {
	clientTransportCredentials, err := credentials.NewClientTLSFromFile(viper.GetString("cert_path"), "")
	conn, err := grpc.DialContext(ctx, "unix://"+socket, grpc.WithTransportCredentials(clientTransportCredentials))
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()

		err := conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("couldn't close grpc connection")
		}
	}()

	return pb.NewAccountsClient(conn), nil
}
