package auth

import (
	authservice "authService/protos/gen/go/sso"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	authservice.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authservice.RegisterAuthServer(gRPC, &ServerAPI{})
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *authservice.LoginRequest,
) (*authservice.LoginResponse, error) {
	return &authservice.LoginResponse{
		AccessToken:  fmt.Sprintf(req.Password, req.GetPassword()),
		RefreshToken: "some refresh token",
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *authservice.RegisterRequest) (
	*authservice.RegisterResponse, error) {
	panic("implement me")
}

func (s *ServerAPI) GetRole(ctx context.Context, req *authservice.GetRoleRequest) (
	*authservice.GetRoleResponse, error) {
	panic("implement me")
}
