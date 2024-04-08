package sqlite

import (
	"authService/internal/domain/models"
	"authService/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) App(ctx context.Context, appID uint32) (models.App, error) {
	//TODO implement me
	panic("implement me")
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, username string, passHash []byte) (uint, error) {
	const op = "storage.sqlite.SaveUser"
	const defaultRole = 1

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, username, passHash, defaultRole)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, storage.ErrUserExists
		}
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return uint(id), nil
}

func (s *Storage) User(ctx context.Context, username string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.PrepareContext(
		ctx, "SELECT id, username, password_hash, role FROM users WHERE username = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s, %w", op, err)
	}

	var user models.User
	row := stmt.QueryRowContext(ctx, username)

	err = row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s, %w", op, err)
	}

	return user, nil
}

func (s *Storage) Role(ctx context.Context, username string) (uint, error) {
	const op = "storage.sqlite.Role"

	stmt, err := s.db.PrepareContext(ctx, "SELECT role FROM users WHERE username = ?")
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, username)

	var role uint
	err = row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrUserNotFound
		}
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return role, nil
}
