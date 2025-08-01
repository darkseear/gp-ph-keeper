package proto

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	pb "github.com/darkseear/gophkeeper/proto"
	"github.com/darkseear/gophkeeper/server/internal/config"
	"github.com/darkseear/gophkeeper/server/internal/logger"
	"github.com/darkseear/gophkeeper/server/internal/model"
	"github.com/darkseear/gophkeeper/server/internal/storage"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// GophkeeperGRPCServer - структура grpc сервера.
type GophkeeperGRPCServer struct {
	pb.UnimplementedGophkeeperServer
	store storage.StorageInterface
	cfg   *config.Config
}

// NewGophkeeperGRPCServer - экземпляр сервера.
func NewGophkeeperGRPCServer(stor storage.StorageInterface, cfg *config.Config) *GophkeeperGRPCServer {
	return &GophkeeperGRPCServer{
		cfg:   cfg,
		store: stor,
	}
}

// Register - метод регистрации нового пользователя .
func (s *GophkeeperGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Проверка существования пользователя
	if _, err := s.store.GetUserByLogin(ctx, req.Login); err == nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Генерация хеша пароля (соль не нужна)
	hash := hashPassword(req.Password)

	// Создание пользователя
	user := model.User{
		Login:        req.Login,
		PasswordHash: hash,
	}
	if err := s.store.CreateUser(ctx, &user); err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.RegisterResponse{UserId: user.ID.String()}, nil
}

// Login - метод авторизации пользователей.
func (s *GophkeeperGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Получение пользователя
	user, err := s.store.GetUserByLogin(ctx, req.Login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	// Проверка пароля

	if err := compareHash(user.PasswordHash, req.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	userID := user.ID.String()

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	bytesrez := []byte(s.cfg.SecretKey)
	tokenString, err := token.SignedString(bytesrez)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{Token: tokenString, UserId: userID}, nil
}

// Sync - метод синхронизации данных приложений и сервера.
func (s *GophkeeperGRPCServer) Sync(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
	// Аутентификация пользователя
	userID, err := s.authenticate(req.Token)
	if err != nil {
		return nil, err
	}

	// Преобразование UUID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	// Разрешение конфликтов
	serverSecrets, err := s.store.GetSecrets(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get secrets")
	}

	// Определение последних версий
	serverVersions := make(map[string]int32)
	for _, sec := range serverSecrets {
		serverVersions[sec.ID.String()] = sec.Version
	}

	// Применение изменений
	for _, local := range req.LocalSecrets {
		localID, err := uuid.Parse(local.Id)
		if err != nil {
			continue
		}

		// Проверка версии
		if serverVer, exists := serverVersions[local.Id]; exists && local.Version <= serverVer {
			continue
		}

		// Обновление или создание
		secret := model.Secrets{
			ID:       localID,
			UserID:   uid,
			Type:     local.Type,
			Metadata: local.Metadata,
			Data:     local.Data,
			Version:  local.Version,
		}
		if err := s.store.UpsertSecret(ctx, &secret); err != nil {
			logger.Log.Error("Failed to upsert secret", zap.Error(err))
		}
	}

	// Возврат актуальных данных
	updatedSecrets, err := s.store.GetSecrets(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get updated secrets")
	}

	pbSecrets := make([]*pb.Secret, 0, len(updatedSecrets))
	for _, sec := range updatedSecrets {
		pbSecrets = append(pbSecrets, &pb.Secret{
			Id:        sec.ID.String(),
			Type:      sec.Type,
			Metadata:  sec.Metadata,
			Data:      sec.Data,
			Version:   sec.Version,
			UpdatedAt: timestamppb.New(sec.UpdatedAt),
		})
	}

	return &pb.SyncResponse{ServerSecrets: pbSecrets}, nil
}

// GetSecret - метод передлачи приватнеых данных пользователей.
func (s *GophkeeperGRPCServer) GetSecret(ctx context.Context, req *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	userID, err := s.authenticate(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}
	secret, err := s.store.GetSecretById(ctx, uid, req.SecretId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "secret not found")
	}

	return &pb.GetSecretResponse{
		Secret: &pb.Secret{
			Id:        secret.ID.String(),
			Type:      secret.Type,
			Metadata:  secret.Metadata,
			Data:      secret.Data,
			Version:   secret.Version,
			UpdatedAt: timestamppb.New(secret.UpdatedAt),
		},
	}, nil
}

func hashPassword(password string) string {
	pass := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("error hashed password: %w", err)
	}
	return string(hashedPassword)
}

func compareHash(hashed string, password string) error {
	h := []byte(hashed)
	err := bcrypt.CompareHashAndPassword(h, []byte(password))
	return err
}

// authenticate - вспомогатльная функция для проверки токена авторизации.
func (s *GophkeeperGRPCServer) authenticate(token string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		bytesrez := []byte(s.cfg.SecretKey)
		return bytesrez, nil
	})
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "invalid token claims")
	}
	return sub, nil
}
