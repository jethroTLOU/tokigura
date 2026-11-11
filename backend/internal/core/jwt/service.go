package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
	"uuid"

	goJwt "github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret            string
	AccessTokenDuration  int
	RefreshTokenDuration int
}

type Service struct {
	secret               []byte
	accessTokenDuration  int
	refreshTokenDuration int
}

func NewService(cfg Config) *Service {
	return &Service{
		secret:               []byte(cfg.JWTSecret),
		accessTokenDuration:  cfg.AccessTokenDuration,
		refreshTokenDuration: cfg.RefreshTokenDuration,
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
		ExpiresAt: goJwt.NewNumericDate(
			time.Now().Add(time.Duration(s.accessTokenDuration) * time.Minute),
		),
		IssuedAt: goJwt.NewNumericDate(time.Now()),
	}

	token := goJwt.NewWithClaims(goJwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) Parse(tokenString string) (*Claims, error) {
	token, err := goJwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *goJwt.Token) (any, error) {
			if _, ok := token.Method.(*goJwt.SigningMethodHMAC); !ok {
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
