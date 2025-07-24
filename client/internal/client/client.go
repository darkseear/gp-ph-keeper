package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/darkseear/gophkeeper/client/internal/crypto"
	"github.com/darkseear/gophkeeper/client/internal/storage"
	pb "github.com/darkseear/gophkeeper/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NetworkClient - отвечает за сетевое взаимодействие с сервером.
type NetworkClient struct {
	conn       *grpc.ClientConn
	grpcClient pb.GophkeeperClient
}

func NewNetworkClient(serverAddress string) (*NetworkClient, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return &NetworkClient{
		conn:       conn,
		grpcClient: pb.NewGophkeeperClient(conn),
	}, nil
}

// Close - закрываем соединение.
func (n *NetworkClient) Close() error {
	return n.conn.Close()
}

// GophKeeperClient - агрегирует работу с сетью, криптографией и локальным хранилищем.
type GophKeeperClient struct {
	network      *NetworkClient
	localStorage *storage.LocalStorage
	crypto       *crypto.CryptoService
	token        string
	userID       string
}

// NewGophkeeperClient - экземпляр с инициализацией хранилища, криптографии и сетевого клиента.
func NewGophkeeperClient(network *NetworkClient, local *storage.LocalStorage, masterPassword string) (*GophKeeperClient, error) {

	token, userID, _, err := local.GetSession()
	if err == nil && token != "" && userID != "" {
		log.Println("Session found, using saved token and userID")
	} else {
		log.Println("Not found")
	}

	return &GophKeeperClient{
		network:      network,
		localStorage: local,
		crypto:       crypto.New(masterPassword),
		token:        token,
		userID:       userID,
	}, nil
}

// Register - метод регистрации новых юзеров приложения.
func (c *GophKeeperClient) Register(ctx context.Context, login, password string) error {
	resp, err := c.network.grpcClient.Register(ctx, &pb.RegisterRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}

	c.userID = resp.UserId
	return nil
}

// Login - метод авторизации юзеров.
func (c *GophKeeperClient) Login(ctx context.Context, login, password string) error {
	resp, err := c.network.grpcClient.Login(ctx, &pb.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return err
	}

	c.token = resp.Token
	c.userID = resp.UserId
	// Сохраняем токен локально
	if err := c.localStorage.SaveSession(c.token, c.userID, time.Now().Add(24*time.Hour)); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	return nil
}

// AddSecret - метод добавления приватной информации юзеров (текст , карта и.д.).
func (c *GophKeeperClient) AddSecret(ctx context.Context, secretType, description string, data []byte) error {
	// Шифруем данные перед отправкой
	encryptedData, err := c.crypto.Encrypt(data)
	if err != nil {
		return err
	}

	// Создаем метаданные
	metadata := map[string]string{
		"description": description,
		"created_at":  time.Now().Format(time.RFC3339),
	}

	// Добавляем тип в метаданные
	metadata["type"] = secretType

	// Локальное сохранение
	localSecret := &storage.Secret{
		ID:        generateID(),
		Type:      secretType,
		Metadata:  metadata,
		Data:      encryptedData,
		Version:   1,
		UpdatedAt: time.Now(),
	}

	if err := c.localStorage.SaveSecret(localSecret); err != nil {
		return err
	}

	// Автоматическая синхронизация
	return c.Sync(ctx)
}

// GetSecret - метод получения приватной информации юзеров.
func (c *GophKeeperClient) GetSecret(ctx context.Context, secretID string) (*storage.Secret, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Сначала проверяем локальное хранилище
	localSecret, err := c.localStorage.GetSecret(secretID)
	if err == nil {
		// Дешифруем данные
		decrypted, err := c.crypto.Decrypt(localSecret.Data)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
		localSecret.Data = decrypted
		return localSecret, nil
	}

	// Если нет локально, запрашиваем с сервера
	resp, err := c.network.grpcClient.GetSecret(ctx, &pb.GetSecretRequest{
		Token:    c.token,
		SecretId: secretID,
	})
	if err != nil {
		return nil, fmt.Errorf("server error: %w", err)
	}

	// Дешифруем полученные данные
	decrypted, err := c.crypto.Decrypt(resp.Secret.Data)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	// Сохраняем в локальное хранилище
	secret := &storage.Secret{
		ID:        resp.Secret.Id,
		Type:      resp.Secret.Type,
		Metadata:  resp.Secret.Metadata,
		Data:      decrypted,
		Version:   resp.Secret.Version,
		UpdatedAt: resp.Secret.UpdatedAt.AsTime(),
	}

	if err := c.localStorage.SaveSecret(secret); err != nil {
		log.Printf("Warning: failed to cache secret locally: %v", err)
	}

	return secret, nil
}

// Sync - метод синхронизации данных.
func (c *GophKeeperClient) Sync(ctx context.Context) error {
	// Получаем локальные секреты
	localSecrets, err := c.localStorage.GetAllSecrets()
	if err != nil {
		return err
	}

	// Конвертируем в protobuf формат
	var pbSecrets []*pb.Secret
	for _, s := range localSecrets {
		pbSecrets = append(pbSecrets, &pb.Secret{
			Id:        s.ID,
			Type:      s.Type,
			Metadata:  s.Metadata,
			Data:      s.Data,
			Version:   s.Version,
			UpdatedAt: timestamppb.New(s.UpdatedAt),
		})
	}

	// Выполняем синхронизацию
	resp, err := c.network.grpcClient.Sync(ctx, &pb.SyncRequest{
		Token:        c.token,
		LocalSecrets: pbSecrets,
	})
	if err != nil {
		return err
	}

	// Сохраняем полученные данные
	for _, s := range resp.ServerSecrets {
		if err := c.localStorage.SaveSecret(&storage.Secret{
			ID:        s.Id,
			Type:      s.Type,
			Metadata:  s.Metadata,
			Data:      s.Data,
			Version:   s.Version,
			UpdatedAt: s.UpdatedAt.AsTime(),
		}); err != nil {
			log.Printf("Failed to save secret %s: %v", s.Id, err)
		}
	}

	return nil
}

// ensureAuthenticated - метод проверки сессии приложения.
func (c *GophKeeperClient) ensureAuthenticated(ctx context.Context) error {
	if c.token == "" {
		// Пытаемся загрузить сохраненный токен
		token, userID, expires, err := c.localStorage.GetSession()
		if err != nil || time.Now().After(expires) {
			return errors.New("not authenticated, please login first")
		}
		c.token = token
		c.userID = userID
	}
	return nil
}

// Close - метод закрывающий хранилище и соединения.
func (c *GophKeeperClient) Close() error {
	if err := c.localStorage.Close(); err != nil {
		log.Printf("Failed to close storage: %v", err)
	}
	return c.network.conn.Close()
}

// generateID - вспомогательная функция генерации ID для секретов.
func generateID() string {
	return uuid.New().String()
}
