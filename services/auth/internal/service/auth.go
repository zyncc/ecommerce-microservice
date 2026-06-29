package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/IBM/sarama"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/types/dto"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/config"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/repository"
	"github.com/zyncc/ecommerce-microservice/services/auth/internal/repository/models"
	"go.uber.org/zap"
)

type AuthService struct {
	logger        *zap.Logger
	userRepo      *repository.UserRepository
	kafkaProducer sarama.SyncProducer
	env           *config.EnvConfig
}

func NewAuthService(logger *zap.Logger, userRepo *repository.UserRepository, kafkaProducer sarama.SyncProducer, env *config.EnvConfig) *AuthService {
	return &AuthService{
		logger,
		userRepo,
		kafkaProducer,
		env,
	}
}

func (s *AuthService) SignUp(ctx context.Context, req dto.SignUpRequest) (*uuid.UUID, error) {
	// check if user exists
	_, err := s.userRepo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, errors.New("failed to fetch user")
	}

	// hash password
	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		s.logger.Error("failed to hash password", zap.Error(err))
		return nil, errors.New("failed to hash password")
	}

	id := uuid.New()
	if err := s.userRepo.CreateUser(ctx, &models.CreateUserParams{
		ID:             id,
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
		Role:           "user",
	}); err != nil {
		return nil, errors.New("failed to create user")
	}

	// user.signed-up event
	jsonMsg, err := json.Marshal(map[string]string{
		"id":    id.String(),
		"email": req.Email,
	})
	if err != nil {
		s.logger.Error("failed to create kafka json message", zap.String("topic", "user.signed-up"), zap.Error(err))
	}

	msg := &sarama.ProducerMessage{
		Topic: "user.signed-up",
		Key:   sarama.StringEncoder(id.String()),
		Value: sarama.ByteEncoder(jsonMsg),
	}

	_, _, err = s.kafkaProducer.SendMessage(msg)
	if err != nil {
		s.logger.Error("failed to send kafka message", zap.String("topic", "user.signed-up"), zap.Error(err))
	}

	return &id, nil
}

func (s *AuthService) SignIn(ctx context.Context, req dto.SignInRequest) (string, string, error) {
	// check if user exists
	user, err := s.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", "", errors.New("invalid email or password")
		}
		return "", "", errors.New("failed to fetch user")
	}

	// compare hash password
	match, err := argon2id.ComparePasswordAndHash(req.Password, user.HashedPassword)
	if err != nil {
		s.logger.Error("failed to compare hash password", zap.Error(err))
		return "", "", errors.New("failed to compare hash password")
	}

	if !match {
		return "", "", errors.New("invalid credentials")
	}

	// create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        user.ID.String(),
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"exp":        time.Now().Add(30 * time.Minute).Unix(),
		"iat":        time.Now().Unix(),
	})

	tokenString, err := accessToken.SignedString([]byte(s.env.JWTSecret))
	if err != nil {
		return "", "", errors.New("failed to sign token")
	}

	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	refreshToken := base64.RawURLEncoding.EncodeToString(buf)
	hashedRefreshToken := sha256.Sum256([]byte(refreshToken))
	hashString := hex.EncodeToString(hashedRefreshToken[:])

	err = s.userRepo.CreateRefreshToken(ctx, user.ID, hashString)
	if err != nil {
		return "", "", errors.New("failed to create refresh token")
	}

	return tokenString, hashString, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	findRefreshToken, err := s.userRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	findUser, err := s.userRepo.FindUserByID(ctx, findRefreshToken.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        findUser.ID.String(),
		"name":       findUser.Name,
		"email":      findUser.Email,
		"role":       findUser.Role,
		"created_at": findUser.CreatedAt,
		"exp":        time.Now().Add(30 * time.Minute).Unix(),
		"iat":        time.Now().Unix(),
	})

	tokenString, err := accessToken.SignedString([]byte(s.env.JWTSecret))
	if err != nil {
		return "", errors.New("failed to sign token")
	}

	return tokenString, nil
}
