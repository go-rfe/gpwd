package server

import (
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
	"github.com/go-rfe/gpwd/internal/storage/cloud"
)

type Cfg struct {
	ServerAddress string        `mapstructure:"server_address"`
	TokenLifespan time.Duration `mapstructure:"token_lifespan"`
	DatabaseDSN   string        `mapstructure:"database_dsn"`
	CertPath      string        `mapstructure:"server_cert_path"`
	KeyPath       string        `mapstructure:"server_key_path"`
}

type server struct {
	cfg            *Cfg
	secretKey      []byte
	accountStorage cloud.Accounts
	secretsStorage cloud.Secrets
	pb.UnimplementedLoginServer
	pb.UnimplementedSyncServer
}

func NewServer(cfg *Cfg) *server {
	return &server{
		cfg: cfg,
	}
}

func (s *server) Run() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	storage, err := cloud.NewDB(s.cfg.DatabaseDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create storage")
	}
	defer func(closer io.Closer) {
		if err := storage.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close storage")
		}
	}(storage)

	s.accountStorage = storage
	s.secretsStorage = storage

	s.secretKey, err = os.ReadFile(s.cfg.KeyPath)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to read TLS key")
	}

	listener, err := s.createListener()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create listener")
	}

	if err := s.listenEndServe(ctx, listener); err != nil {
		log.Fatal().Err(err).Msg("couldn't start server")
	}
}

func (s *server) createListener() (net.Listener, error) {
	l, err := net.Listen("tcp", s.cfg.ServerAddress)
	if err != nil {
		return nil, err
	}

	return l, err
}

func (s *server) listenEndServe(ctx context.Context, listener net.Listener) error {
	serverTransportCreds, err := credentials.NewServerTLSFromFile(s.cfg.CertPath, s.cfg.KeyPath)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(serverTransportCreds),
		grpc.UnaryInterceptor(s.authUnaryInterceptor),
		grpc.StreamInterceptor(s.authStreamInterceptor),
	)

	pb.RegisterLoginServer(grpcServer, s)
	pb.RegisterSyncServer(grpcServer, s)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(listener)
}
