package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jethroTLOU/tokigura/internal/config"
)

const (
	accessTokenDuration = 15
	refreshTokenBytes   = 32
)

type Service struct {
	secret []byte
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		secret: []byte(cfg.JWTSecret),
	}
}

func (s *Service) GenerateTokenPair(userID uuid.UUID, role string) (*TokenPair, error) {
	refresh, err := s.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	access, err := s.GenerateAccessToken(userID, role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{RefreshToken: refresh, AccessToken: access}, nil
}

func (s *Service) GenerateRefreshToken() (string, error) {
	b := make([]byte, refreshTokenBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Service) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}
