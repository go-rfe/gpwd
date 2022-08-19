package agent

import (
	"context"

	"github.com/google/uuid"

	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
	"github.com/go-rfe/gpwd/internal/storage/local"
)

// CreateSecret creates secret
func (a *agent) CreateSecret(ctx context.Context, request *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	secret := request.Secret

	secret.ID = uuid.New().String()

	log.Info().Msgf("CreateSecret secret %s", secret.ID)

	var err error
	secret.Data, err = a.encrypt(secret.GetData())
	if err != nil {
		return nil, err
	}

	id, err := a.secretsStorage.CreateSecret(ctx, secret)
	if err != nil {
		return nil, err
	}

	return &pb.CreateSecretResponse{
		Error: "",
		Id:    id,
	}, nil
}

// ListSecrets list secrets
func (a *agent) ListSecrets(ctx context.Context, _ *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	secrets, err := a.secretsStorage.ListSecrets(ctx)
	if err != nil {
		return nil, err
	}

	// Filter deleted secrets
	actualSecrets := make([]*pb.Secret, 0, len(secrets))
	for _, secret := range secrets {
		if !secret.Status.Deleted {
			actualSecrets = append(actualSecrets, secret)
		}
	}

	return &pb.ListSecretsResponse{
		Secrets: actualSecrets,
	}, nil
}

// GetSecret returns secret
func (a *agent) GetSecret(ctx context.Context, request *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	id := request.GetId()

	secret, err := a.secretsStorage.GetSecret(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetSecretResponse{
		Error:  "",
		Secret: secret,
	}, nil
}

// UpdateSecret updates secret
func (a *agent) UpdateSecret(ctx context.Context, request *pb.UpdateSecretRequest) (*pb.UpdateSecretResponse, error) {
	secret := request.Secret

	log.Info().Msgf("UpdateSecret secret %s", secret.ID)

	var err error
	secret.Data, err = a.encrypt(secret.GetData())
	if err != nil {
		return nil, err
	}

	if err := a.secretsStorage.UpdateSecret(ctx, secret); err != nil {
		return nil, err
	}

	return &pb.UpdateSecretResponse{
		Error: "",
		Id:    secret.GetID(),
	}, nil
}

// DeleteSecret deletes secret
func (a *agent) DeleteSecret(ctx context.Context, request *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	secret := request.GetSecret()

	err := a.secretsStorage.DeleteSecret(ctx, secret)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteSecretResponse{
		Error: "",
	}, nil
}

// CreateAccount creates account
func (a *agent) CreateAccount(ctx context.Context, request *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	// check first if account exist
	existingAccount, err := a.accountsStorage.GetAccount(ctx)
	if err != nil && existingAccount == nil {
		return nil, err
	}
	if existingAccount != nil {
		return nil, local.ErrAccountExists
	}

	account := request.Account

	account.ID = uuid.New().String()

	log.Info().Msgf("Create account %s", account.ID)

	account.UserPassword, err = a.encrypt(account.GetUserPassword())
	if err != nil {
		return nil, err
	}

	id, err := a.accountsStorage.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return &pb.CreateAccountResponse{
		Error: "",
		Id:    id,
	}, nil
}

// GetAccount returns account
func (a *agent) GetAccount(ctx context.Context, _ *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := a.accountsStorage.GetAccount(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Error:   "",
		Account: account,
	}, nil
}

// UpdateAccount updates account
func (a *agent) UpdateAccount(ctx context.Context, request *pb.UpdateAccountRequest) (*pb.UpdateAccountResponse, error) {
	account := request.Account

	log.Info().Msgf("Update account %s", account.ID)

	// check first if account exist
	existingAccount, err := a.accountsStorage.GetAccount(ctx)
	if err != nil && existingAccount == nil {
		return nil, err
	}

	if account.GetServerAddress() != "" {
		existingAccount.ServerAddress = account.GetServerAddress()
	}

	if account.GetUserName() != "" {
		existingAccount.UserName = account.GetUserName()
	}

	if account.GetUserPassword() != nil {
		existingAccount.UserPassword, err = a.encrypt(account.GetUserPassword())
		if err != nil {
			return nil, err
		}
	}

	if err := a.accountsStorage.UpdateAccount(ctx, account); err != nil {
		return nil, err
	}

	return &pb.UpdateAccountResponse{
		Error: "",
	}, nil
}

// DeleteAccount deletes account
func (a *agent) DeleteAccount(ctx context.Context, _ *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	log.Info().Msgf("Delete existing account")

	err := a.accountsStorage.DeleteAccount(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteAccountResponse{
		Error: "",
	}, nil
}
