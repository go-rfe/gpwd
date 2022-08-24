package server

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/go-rfe/gpwd/internal/logging/log"
	pb "github.com/go-rfe/gpwd/internal/proto"
)

func (s *server) GenerateJWT(auth *pb.Auth) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   auth.GetUsername(),
		ExpiresAt: time.Now().Add(s.cfg.TokenLifespan).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *server) VerifyJWT(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return s.secretKey, nil
		},
	)

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.Subject, nil
}

func (s *server) authUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if err := s.authorize(ctx, info.FullMethod); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func (s *server) authStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if err := s.authorize(stream.Context(), info.FullMethod); err != nil {
		return err
	}

	return handler(srv, stream)
}

func (s *server) authorize(ctx context.Context, method string) error {
	log.Debug().Msg(method)

	if method == "/proto.Login/RegisterAccount" || method == "/proto.Login/Login" {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["jwt"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	_, err := s.VerifyJWT(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return nil
}

func (s *server) mustReturnUsernameFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Fatal().Msg("BUG: cannot get username from context")
	}

	values := md["jwt"]
	if len(values) == 0 {
		log.Fatal().Msg("BUG: cannot get username from context")
	}

	accessToken := values[0]
	username, err := s.VerifyJWT(accessToken)
	if err != nil {
		log.Fatal().Msg("BUG: cannot get username from context")
	}

	return username
}
