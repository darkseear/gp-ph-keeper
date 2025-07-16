package storage

import (
	"context"

	"github.com/darkseear/gophkeeper/server/internal/model"
	"github.com/google/uuid"
)

// StorageInterface - интерфейс методов хранилища.
type StorageInterface interface {
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	GetSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secrets, error)
	GetSecretById(ctx context.Context, userID uuid.UUID, secretID string) (*model.Secrets, error)
	UpsertSecret(ctx context.Context, secret *model.Secrets) error
	Close() error
}
