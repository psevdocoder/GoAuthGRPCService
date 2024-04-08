package auth

import (
	"authService/internal/domain/dto"
	"authService/internal/services/auth"
	authservice "authService/protos/gen/go/sso"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, input dto.LoginInput) (string, error)
	Register(ctx context.Context, input dto.RegisterInput) (uint, error)
	Role(ctx context.Context, username string) (uint, error)
}

type ServerAPI struct {
	authservice.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authservice.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *authservice.LoginRequest,
) (*authservice.LoginResponse, error) {

	input := dto.LoginInput{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	}

	if err := input.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accessToken, err := s.auth.Login(ctx, input)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authservice.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: "unimplemented", // TODO: refresh token
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *authservice.RegisterRequest) (
	*authservice.RegisterResponse, error) {

	input := dto.RegisterInput{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}

	if err := input.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.auth.Register(ctx, input)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authservice.RegisterResponse{UserId: uint32(id)}, nil
}

func (s *ServerAPI) GetRole(ctx context.Context, req *authservice.GetRoleRequest) (
	*authservice.GetRoleResponse, error) {
	panic("implement me")
}
