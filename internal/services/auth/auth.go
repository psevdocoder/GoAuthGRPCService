package auth

import (
	"authService/internal/domain/dto"
	"authService/internal/domain/models"
	"authService/internal/storage"
	"authService/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type UserSaver interface {
	SaveUser(ctx context.Context, username string, passwordHash []byte) (uint, error)
}

type UserProvider interface {
	User(ctx context.Context, username string) (models.User, error)
	Role(ctx context.Context, userID uint32) (uint, error)
}

type AppProvider interface {
	App(ctx context.Context, appID uint32) (models.App, error)
}

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrAppNotFound        = errors.New("app not found")
)

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Register(ctx context.Context, input dto.RegisterInput) (uint, error) {
	const op = "auth.Register"
	log := a.log.With(slog.String("op", op))

	log.Info("register attempt", slog.String("username", input.Username))

	passHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, input.Username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user already exists", slog.String("err", err.Error()))
			return 0, ErrUserExists
		}
		log.Error("registration failed", slog.String("err", err.Error()))
		return 0, errors.New("registration failed")
	}

	return id, nil
}

func (a *Auth) Login(ctx context.Context, input dto.LoginInput) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))

	log.Info("login attempt", slog.String("username", input.Username))

	user, err := a.usrProvider.User(ctx, input.Username)
	if err != nil {

		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("err", err.Error()))
			return "", ErrInvalidCredentials
		}

		log.Error("failed to get user", slog.String("err", err.Error()))
		return "", errors.New("failed to get user")
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password)); err != nil {
		log.Error("Invalid credentials", slog.String("err", err.Error()))
		return "", ErrInvalidCredentials
	}

	app, err := a.appProvider.App(ctx, input.AppID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("app not found", slog.String("err", err.Error()))
			return "", ErrAppNotFound
		}
		log.Error("failed to get app", slog.String("err", err.Error()))
		return "", errors.New("failed to get app")
	}

	log.Info("login success", slog.String("username", input.Username))

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to create token", slog.String("err", err.Error()))
		return "", errors.New("failed to create token")
	}

	return token, nil
}

func (a *Auth) Role(ctx context.Context, userID uint32) (uint32, error) {
	const op = "auth.Role"
	log := a.log.With(slog.String("op", op))

	log.Info("get role attempt", slog.Uint64("userID", uint64(userID)))

	role, err := a.usrProvider.Role(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.String("err", err.Error()))
			return 0, ErrUserNotFound
		}

		log.Error("failed to get role", slog.String("err", err.Error()))
		return 0, errors.New("failed to get role")
	}

	log.Info("get role success", slog.Uint64("userID", uint64(userID)), slog.Int("role", int(role)))

	return uint32(role), nil
}
