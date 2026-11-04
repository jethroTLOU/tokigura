package jwt

// TODO: много чего делать
import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jethroTLOU/tokigura/internal/core/config"
)

type Service struct {
	secret               []byte
	accessTokenDuration  int
	refreshTokenDuration int
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		secret:               []byte(cfg.Auth.JWTSecret),
		accessTokenDuration:  cfg.Auth.AccessTokeDuration,
		refreshTokenDuration: cfg.Auth.RefreshTokeDuration,
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

	return &TokenPair{
		RefreshToken:     refresh,
		AccessToken:      access,
		RefreshExpiresAt: time.Now().AddDate(0, 0, s.refreshTokenDuration),
	}, nil
}

func (s *Service) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Service) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		Role:   role,
		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(time.Duration(s.accessTokenDuration) * time.Minute),
		),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
