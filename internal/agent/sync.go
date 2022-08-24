package agent

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
	"github.com/go-rfe/gpwd/internal/storage/local"
	"github.com/go-rfe/gpwd/internal/syncer"
)

func (a *agent) syncWorker(ctx context.Context) {
	registerTicker := time.NewTicker(a.cfg.SyncInterval)
	defer registerTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-registerTicker.C:
			syncCtx, cancel := context.WithTimeout(ctx, a.cfg.SyncInterval)

			client, err := a.getSyncClient(syncCtx)
			if err != nil {
				log.Error().Err(err).Msg("a error occurred during sync client creation")
				cancel()
				continue
			}

			if err := a.syncDeleted(ctx, client); err != nil {
				log.Error().Err(err).Msg("a error occurred during sync deleted secrets")
			}
			if err := a.syncUpdated(ctx, client); err != nil {
				log.Error().Err(err).Msg("a error occurred during sync updated secrets")
			}
			if err := a.syncCreated(ctx, client); err != nil {
				log.Error().Err(err).Msg("a error occurred during sync created secrets")
			}
			if err := a.sync(ctx, client); err != nil {
				log.Error().Err(err).Msg("a error occurred during sync secrets")
			}

			cancel()
		}
	}
}

func (a *agent) getSyncClient(ctx context.Context) (pb.SyncClient, error) {
	account, err := a.accountsStorage.GetAccount(ctx)
	if err != nil {
		return nil, err
	}

	registered := account.Registered

	account.UserPassword, err = a.decrypt(account.GetUserPassword())
	if err != nil {
		return nil, err
	}

	client, err := syncer.NewSyncer(ctx, account)
	if err != nil {
		return nil, err
	}

	if account.Registered && !registered {
		if err := a.accountsStorage.UpdateAccount(ctx, account); err != nil {
			return nil, err
		}
	}

	return client, err
}

func (a *agent) syncDeleted(ctx context.Context, client pb.SyncClient) error {
	secrets, err := a.secretsStorage.ListSecrets(ctx)
	if err != nil {
		return err
	}

	stream, err := client.SyncDeleted(ctx)
	for _, secret := range secrets {
		if secret.Status.Deleted && !secret.Status.Synced {
			err = stream.Send(&pb.SyncRequest{Secret: secret})
			if err != nil {
				return err
			}

			secret.Status.Synced = true

			if err := a.secretsStorage.UpdateSecret(ctx, secret); err != nil {
				return err
			}
		}
	}
	recv, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if recv.GetError() != "" {
		return errors.New(recv.GetError())
	}

	return nil
}

func (a *agent) syncUpdated(ctx context.Context, client pb.SyncClient) error {
	secrets, err := a.secretsStorage.ListSecrets(ctx)
	if err != nil {
		return err
	}

	stream, err := client.SyncUpdated(ctx)
	for _, secret := range secrets {
		if !secret.Status.Deleted && !secret.Status.Synced && secret.UpdatedAt != nil {
			err = stream.Send(&pb.SyncRequest{Secret: secret})
			if err != nil {
				return err
			}

			secret.Status.Synced = true

			if err := a.secretsStorage.UpdateSecret(ctx, secret); err != nil {
				return err
			}
		}
	}
	recv, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if recv.GetError() != "" {
		return errors.New(recv.GetError())
	}

	return nil
}

func (a *agent) syncCreated(ctx context.Context, client pb.SyncClient) error {
	secrets, err := a.secretsStorage.ListSecrets(ctx)
	if err != nil {
		return err
	}

	stream, err := client.SyncCreated(ctx)
	for _, secret := range secrets {
		if !secret.Status.Deleted && !secret.Status.Synced && secret.UpdatedAt == nil {
			err = stream.Send(&pb.SyncRequest{Secret: secret})
			if err != nil {
				return err
			}

			secret.Status.Synced = true

			if err := a.secretsStorage.UpdateSecret(ctx, secret); err != nil {
				return err
			}
		}
	}
	recv, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	if recv.GetError() != "" {
		return errors.New(recv.GetError())
	}

	return nil
}

func (a *agent) sync(ctx context.Context, client pb.SyncClient) error {
	var secrets []*pb.Secret

	stream, err := client.Sync(ctx, &pb.SyncRequest{})
	if err != nil {
		return err
	}

	for {
		message, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		secrets = append(secrets, message.GetSecret())
	}

	if err := stream.CloseSend(); err != nil {
		return err
	}

	for _, secret := range secrets {
		localSecret, err := a.secretsStorage.GetSecret(ctx, secret.GetID())
		if err != nil && !errors.Is(err, local.ErrNoSecretFound) {
			return err
		}

		secret.Status.Synced = true

		switch {
		case secret.GetStatus().GetDeleted():
			if err := a.secretsStorage.DeleteSecret(ctx, secret); err != nil {
				return err
			}
		case errors.Is(err, local.ErrNoSecretFound) && !secret.GetStatus().GetDeleted():
			if _, err := a.secretsStorage.CreateSecret(ctx, secret); err != nil {
				return err
			}
		case localSecret.GetUpdatedAt().AsTime().Before(secret.GetUpdatedAt().AsTime()) && !secret.GetStatus().GetDeleted():
			if err := a.secretsStorage.UpdateSecret(ctx, secret); err != nil {
				return err
			}
		}
	}

	return nil
}
