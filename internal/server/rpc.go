package server

import (
	"context"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"

	pb "github.com/go-rfe/gpwd/internal/proto"
)

func (s *server) RegisterAccount(ctx context.Context, request *pb.RegisterAccountRequest) (*pb.RegisterAccountResponse, error) {
	auth := request.GetAuth()

	var err error
	auth.Password, err = bcrypt.GenerateFromPassword(auth.GetPassword(), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = s.accountStorage.CreateAccount(ctx, auth)
	if err != nil {
		return nil, err
	}

	token, err := s.GenerateJWT(auth)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterAccountResponse{
		Error: "",
		Token: token,
	}, nil
}

func (s *server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	auth := request.GetAuth()

	existingAuth, err := s.accountStorage.GetByName(ctx, auth.GetUsername())
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(existingAuth.GetPassword(), auth.GetPassword()); err != nil {
		return nil, err
	}

	token, err := s.GenerateJWT(auth)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Error: "",
		Token: token,
	}, nil
}

func (s *server) SyncDeleted(stream pb.Sync_SyncDeletedServer) error {
	var secrets []*pb.Secret
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

	username := s.mustReturnUsernameFromContext(stream.Context())

	auth := &pb.Auth{
		Username: username,
	}

	err := s.secretsStorage.DeleteSecrets(stream.Context(), auth, secrets)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&pb.SyncResponse{})
}

func (s *server) SyncUpdated(stream pb.Sync_SyncUpdatedServer) error {
	var secrets []*pb.Secret
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

	username := s.mustReturnUsernameFromContext(stream.Context())

	auth := &pb.Auth{
		Username: username,
	}

	err := s.secretsStorage.UpdateSecrets(stream.Context(), auth, secrets)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&pb.SyncResponse{})
}

func (s *server) SyncCreated(stream pb.Sync_SyncCreatedServer) error {
	var secrets []*pb.Secret
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

	username := s.mustReturnUsernameFromContext(stream.Context())

	auth := &pb.Auth{
		Username: username,
	}

	err := s.secretsStorage.CreateSecrets(stream.Context(), auth, secrets)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&pb.SyncResponse{})
}

func (s *server) Sync(_ *pb.SyncRequest, stream pb.Sync_SyncServer) error {
	username := s.mustReturnUsernameFromContext(stream.Context())
	auth := &pb.Auth{
		Username: username,
	}

	secrets, err := s.secretsStorage.ListSecrets(stream.Context(), auth)
	if err != nil {
		return err
	}

	for _, secret := range secrets {
		if err := stream.Send(&pb.SyncResponse{
			Secret: secret,
		}); err != nil {
			return err
		}
	}

	return nil
}
