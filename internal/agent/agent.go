package agent

import (
	"context"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-rfe/gpwd/internal/encryption"
	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto" // import protobufs
	"github.com/go-rfe/gpwd/internal/storage/local"
)

type Cfg struct {
	SocketPath     string        `mapstructure:"socket_path"`
	SyncInterval   time.Duration `mapstructure:"sync_interval"`
	StorePath      string        `mapstructure:"store_path"`
	CertPath       string        `mapstructure:"agent_cert_path"`
	KeyPath        string        `mapstructure:"agent_key_path"`
	MasterPassword []byte        `mapstructure:"master_password"`
}

type agent struct {
	cfg             *Cfg
	secretsStorage  local.Secrets
	accountsStorage local.Accounts
	mu              sync.RWMutex
	encrypt         func([]byte) ([]byte, error)
	decrypt         func([]byte) ([]byte, error)
	pb.UnimplementedSecretsServer
	pb.UnimplementedAccountsServer
}

func NewAgent(cfg *Cfg) *agent {
	return &agent{
		cfg: cfg,
	}
}

func (a *agent) Run() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	if err := a.createDirs(); err != nil {
		log.Fatal().Err(err).Msg("couldn't create agent working directory")
	}

	storage, err := local.NewSQLiteStorage(a.cfg.StorePath, a.cfg.MasterPassword)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create storage")
	}
	defer func(closer io.Closer) {
		if err := closer.Close(); err != nil {
			log.Error().Err(err).Msgf("failed to close storage")
		}
	}(storage)

	a.secretsStorage = storage
	a.accountsStorage = storage

	encrypt, decrypt, err := encryption.GetCrypto(a.cfg.MasterPassword)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to init encryption")
	}

	a.encrypt = encrypt
	a.decrypt = decrypt

	listener, err := a.createListener()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create listener")
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.listenEndServe(ctx, listener); err != nil {
			log.Fatal().Err(err).Msg("couldn't start agent")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.syncWorker(ctx)
	}()

	wg.Wait()
}

func (a *agent) createDirs() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	for _, path := range []string{a.cfg.SocketPath, a.cfg.StorePath} {
		if strings.HasPrefix(path, "~/") {
			a.cfg.SocketPath = filepath.Join(home, path[1:])
		}

		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
	}

	return nil
}

func (a *agent) createListener() (net.Listener, error) {
	l, err := net.Listen("unix", a.cfg.SocketPath)
	if err != nil {
		return nil, err
	}

	return l, err
}

func (a *agent) listenEndServe(ctx context.Context, listener net.Listener) error {
	serverTransportCreds, err := credentials.NewServerTLSFromFile(a.cfg.CertPath, a.cfg.KeyPath)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(grpc.Creds(serverTransportCreds))

	pb.RegisterSecretsServer(grpcServer, a)
	pb.RegisterAccountsServer(grpcServer, a)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(listener)
}
