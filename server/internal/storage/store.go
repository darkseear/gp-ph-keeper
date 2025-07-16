package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/darkseear/gophkeeper/server/internal/config"
	"github.com/darkseear/gophkeeper/server/internal/logger"
	"github.com/darkseear/gophkeeper/server/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotFound = errors.New("not found")

// Store - структура хранилища сервера.
type Store struct {
	db  *sqlx.DB
	cfg *config.Config
}

// NewStore экземпляр хранилища.
func NewStore(cfg *config.Config) (*Store, error) {
	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN)
	if err != nil {
		logger.Log.Error("Error create storage DB", zap.Error(err))
		return nil, err
	}
	//Проверка соединения
	if err := db.Ping(); err != nil {
		logger.Log.Error("Ping failed", zap.Error(err))
		return nil, err
	}
	// Устанавливаем настройки пула
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := ApplyMigrations(db); err != nil {
		logger.Log.Error("Migrations failed", zap.Error(err))
		return nil, err
	}

	return &Store{
		db:  db,
		cfg: cfg,
	}, nil
}

// ApplyMigrations - выполняет миграции базы данных.
func ApplyMigrations(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			salt TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS secrets (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			type TEXT NOT NULL,
			metadata JSONB,
			data BYTEA,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS secrets_user_id_idx ON secrets(user_id);
		CREATE INDEX IF NOT EXISTS secrets_updated_at_idx ON secrets(updated_at);
	`)
	return err
}

// User методы
// GetUserByLogin - получение юзера по логину.
func (s *Store) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	const query = `SELECT id, login, password_hash, salt, created_at FROM users WHERE login = $1`

	var user model.User
	err := s.db.GetContext(ctx, &user, query, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "Not found User")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// CreateUser - метод созхдания нового пользователя.
func (s *Store) CreateUser(ctx context.Context, user *model.User) error {
	const query = `
		INSERT INTO users (id, login, password_hash, salt, created_at)
		VALUES (:id, :login, :password_hash, :salt, :created_at)
	`

	user.ID = uuid.New()
	user.CreatedAt = time.Now()

	_, err := s.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetSecrets - получение приватных данных владельцем.
func (s *Store) GetSecrets(ctx context.Context, userID uuid.UUID) ([]*model.Secrets, error) {
	const query = `
		SELECT id, user_id, type, metadata, data, version, created_at, updated_at
		FROM secrets
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	var secrets []*model.Secrets
	rows, err := s.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query secrets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sec model.Secrets
		var meta []byte

		if err := rows.Scan(
			&sec.ID,
			&sec.UserID,
			&sec.Type,
			&meta,
			&sec.Data,
			&sec.Version,
			&sec.CreatedAt,
			&sec.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}

		// Декодируем метаданные из JSON
		if err := json.Unmarshal(meta, &sec.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		secrets = append(secrets, &sec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return secrets, nil
}

// GetSecret - получение конкретного секрета.
func (s *Store) GetSecretById(ctx context.Context, userID uuid.UUID, secretID string) (*model.Secrets, error) {
	const query = `
		SELECT id, user_id, type, metadata, data, version, created_at, updated_at 
		FROM secrets 
		WHERE id = $1 AND user_id = $2
	`

	var (
		sec      model.Secrets
		metaJSON []byte
	)

	// Проверяем что secretID - валидный UUID
	if _, err := uuid.Parse(secretID); err != nil {
		return nil, fmt.Errorf("invalid secret ID format: %w", err)
	}

	row := s.db.QueryRowContext(ctx, query, secretID, userID)

	err := row.Scan(
		&sec.ID,
		&sec.UserID,
		&sec.Type,
		&metaJSON,
		&sec.Data,
		&sec.Version,
		&sec.CreatedAt,
		&sec.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query secret: %w", err)
	}

	// Декодируем метаданные из JSON
	if err := json.Unmarshal(metaJSON, &sec.Metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	// Дополнительная проверка принадлежности пользователю
	if sec.UserID != userID {
		return nil, fmt.Errorf("err access denaid: %w", err)
	}

	return &sec, nil
}

// UpsertSecret - обновление приватных данных.
func (s *Store) UpsertSecret(ctx context.Context, secret *model.Secrets) error {
	const query = `
		INSERT INTO secrets (id, user_id, type, metadata, data, version, updated_at)
		VALUES (:id, :user_id, :type, :metadata, :data, :version, :updated_at)
		ON CONFLICT (id) DO UPDATE SET
			type = EXCLUDED.type,
			metadata = EXCLUDED.metadata,
			data = EXCLUDED.data,
			version = EXCLUDED.version,
			updated_at = EXCLUDED.updated_at
		WHERE secrets.version < EXCLUDED.version
	`

	// Подготовка метаданных
	meta, err := json.Marshal(secret.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Увеличиваем версию при обновлении
	secret.Version++
	secret.UpdatedAt = time.Now()

	_, err = s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":         secret.ID,
		"user_id":    secret.UserID,
		"type":       secret.Type,
		"metadata":   meta,
		"data":       secret.Data,
		"version":    secret.Version,
		"updated_at": secret.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert secret: %w", err)
	}

	return nil
}

// Close - для закрывания.
func (s *Store) Close() error {
	return s.db.Close()

}
