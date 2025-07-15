package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Secret - структура приватных данных.
type Secret struct {
	ID        string
	Type      string
	Metadata  map[string]string
	Data      []byte
	Version   int32
	UpdatedAt time.Time
}

// LocalStorage - структура локального хранилища.
type LocalStorage struct {
	db *sql.DB
}

// NewLocalStorage - экземпляр с инициализацией бд локального хранилища.
func NewLocalStorage() (*LocalStorage, error) {
	db, err := sql.Open("sqlite3", "gophkeeper.db")
	if err != nil {
		return nil, err
	}

	if err := initDB(db); err != nil {
		return nil, err
	}

	return &LocalStorage{db: db}, nil
}

// initDB - инициализатор бд .
func initDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS secrets (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			metadata TEXT NOT NULL,
			data BLOB NOT NULL,
			version INTEGER NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS sessions (
			token TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_secrets_type ON secrets(type);
		CREATE INDEX IF NOT EXISTS idx_secrets_updated ON secrets(updated_at);
	`)
	return err
}

// SaveSession - сессия посещения.
func (s *LocalStorage) SaveSession(token, userID string, expires time.Time) error {
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?)",
		token, userID, expires,
	)
	return err
}

// GetSession - метод , там  получаем сессию.
func (s *LocalStorage) GetSession() (token, userID string, expires time.Time, err error) {
	row := s.db.QueryRow("SELECT token, user_id, expires_at FROM sessions LIMIT 1")
	err = row.Scan(&token, &userID, &expires)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", time.Time{}, nil
	}
	return
}

// SaveSecret - локальное хранение приватных данных.
func (s *LocalStorage) SaveSecret(secret *Secret) error {
	meta, err := json.Marshal(secret.Metadata)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`INSERT OR REPLACE INTO secrets 
		(id, type, metadata, data, version, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		secret.ID,
		secret.Type,
		string(meta),
		secret.Data,
		secret.Version,
		secret.UpdatedAt,
	)
	return err
}

// GetSecret - получаем приватные данные из хранилища.
func (s *LocalStorage) GetSecret(id string) (*Secret, error) {
	row := s.db.QueryRow(
		"SELECT id, type, metadata, data, version, updated_at FROM secrets WHERE id = ?",
		id,
	)

	var sec Secret
	var metaStr string
	err := row.Scan(
		&sec.ID,
		&sec.Type,
		&metaStr,
		&sec.Data,
		&sec.Version,
		&sec.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(metaStr), &sec.Metadata); err != nil {
		return nil, err
	}

	return &sec, nil
}

// GetSecretsByType - получаем данные по типу.
func (s *LocalStorage) GetSecretsByType(secretType string) ([]*Secret, error) {
	rows, err := s.db.Query(
		"SELECT id, type, metadata, data, version, updated_at FROM secrets WHERE type = ?",
		secretType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*Secret
	for rows.Next() {
		var sec Secret
		var metaStr string

		if err := rows.Scan(
			&sec.ID,
			&sec.Type,
			&metaStr,
			&sec.Data,
			&sec.Version,
			&sec.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(metaStr), &sec.Metadata); err != nil {
			return nil, err
		}

		secrets = append(secrets, &sec)
	}

	return secrets, nil
}

// GetAllSecrets - метод для получения всех данных.
func (s *LocalStorage) GetAllSecrets() ([]*Secret, error) {
	rows, err := s.db.Query(
		"SELECT id, type, metadata, data, version, updated_at FROM secrets",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*Secret
	for rows.Next() {
		var sec Secret
		var metaStr string

		if err := rows.Scan(
			&sec.ID,
			&sec.Type,
			&metaStr,
			&sec.Data,
			&sec.Version,
			&sec.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(metaStr), &sec.Metadata); err != nil {
			return nil, err
		}

		secrets = append(secrets, &sec)
	}

	return secrets, nil
}

// Close - Закрываем хранилище.
func (s *LocalStorage) Close() error {
	return s.db.Close()
}
