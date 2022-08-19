package syncer

import (
	"context"
	"errors"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
)

func NewSyncer(ctx context.Context, account *pb.Account) (pb.SyncClient, error) {
	auth := &pb.Auth{
		Username: account.GetUserName(),
		Password: account.GetUserPassword(),
	}

	login := getGRPCLoginClient(ctx, account.GetServerAddress())

	if !account.Registered {
		resp, err := login.RegisterAccount(ctx, &pb.RegisterAccountRequest{Auth: auth})
		if err != nil {
			return nil, err
		}

		if resp.GetError() != "" {
			return nil, errors.New(resp.GetError())
		}

	}

	account.Registered = true

	token := tokenGenerator(ctx, auth, login)
	interceptor := authInterceptor(token)

	return getGRPCSyncClient(ctx, account.GetServerAddress(), interceptor), nil
}

func getGRPCLoginClient(ctx context.Context, serverAddress string) pb.LoginClient {
	clientTransportCredentials, err := credentials.NewClientTLSFromFile(viper.GetString("server_cert_path"), "")
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(clientTransportCredentials))
	if err != nil {
		log.Fatal().Err(err).Msgf("Couldn't create grpc connection %s", serverAddress)
	}

	go func() {
		<-ctx.Done()

		err := conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("couldn't close grpc connection")
		}
	}()

	return pb.NewLoginClient(conn)
}

func getGRPCSyncClient(ctx context.Context, serverAddress string, interceptor grpc.StreamClientInterceptor) pb.SyncClient {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
	}

	clientTransportCredentials, err := credentials.NewClientTLSFromFile(viper.GetString("server_cert_path"), "")

	conn, err := grpc.DialContext(ctx, serverAddress,
		grpc.WithTransportCredentials(clientTransportCredentials),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(grpc_retry.StreamClientInterceptor(opts...), interceptor)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(grpc_retry.UnaryClientInterceptor(opts...))),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Couldn't create grpc connection %s", serverAddress)
	}

	go func() {
		<-ctx.Done()

		err := conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("couldn't close grpc connection")
		}
	}()

	return pb.NewSyncClient(conn)
}

func authInterceptor(token func() string) grpc.StreamClientInterceptor {
	return func(ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, "jwt", token()), desc, cc, method, opts...)
	}
}

func tokenGenerator(ctx context.Context, auth *pb.Auth, client pb.LoginClient) func() string {
	return func() string {
		login, err := client.Login(ctx, &pb.LoginRequest{Auth: auth})
		if err != nil {
			log.Error().Err(err).Msg("couldn't generate token")
			return ""
		}

		return login.GetToken()
	}
}
