package session

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/netip"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/config"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
	"github.com/jethroTLOU/tokigura/pkg/cryptoutil"
)

type repository interface {
	Create(ctx context.Context, s *Session) error
	DeleteByTokenHash(ctx context.Context, hash string) error
	UpdateRotated(ctx context.Context, params UpdateRotatedParams) error
	GetByTokenHash(ctx context.Context, hash string) (*Session, error)
	DeleteExpired(ctx context.Context) error
}
type jwtService interface {
	GenerateTokenPair(userID uuid.UUID, role user.Role) (*jwt.TokenPair, error)
}

type Service struct {
	repository repository
	jwtService jwtService
	config     *config.Config
}

func NewService(repo repository, jwt jwtService, config *config.Config) *Service {
	return &Service{
		repository: repo,
		jwtService: jwt,
		config:     config,
	}
}

func (s *Service) Create(ctx context.Context, dto CreateSessionDTO) (*jwt.TokenPair, error) {
	tokens, err := s.jwtService.GenerateTokenPair(dto.UserID, dto.Role)
	if err != nil {
		return nil, err
	}

	sess := &Session{
		UserID:        dto.UserID,
		TokenHash:     cryptoutil.SHA256Hex(tokens.RefreshToken),
		DeviceInfo:    dto.UserAgent,
		DeviceHash:    s.generateDeviceHash(dto.DeviceID, dto.UserID),
		IPAddress:     dto.IP,
		LastIPAddress: dto.IP,
		ExpiresAt:     tokens.RefreshExpiresAt,
	}

	if err = s.repository.Create(ctx, sess); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Service) GetSessionByToken(ctx context.Context, refreshToken string) (*Session, error) {
	hash := cryptoutil.SHA256Hex(refreshToken)
	return s.repository.GetByTokenHash(ctx, hash)
}

func (s *Service) DeleteByToken(ctx context.Context, refreshToken string) error {
	hash := cryptoutil.SHA256Hex(refreshToken)
	return s.repository.DeleteByTokenHash(ctx, hash)
}

func (s *Service) RotateSession(
	ctx context.Context,
	sessionID, userID uuid.UUID,
	role user.Role,
	ip netip.Addr,
) (*jwt.TokenPair, error) {
	tokens, err := s.jwtService.GenerateTokenPair(userID, role)
	if err != nil {
		return nil, err
	}

	newHash := cryptoutil.SHA256Hex(tokens.RefreshToken)

	err = s.repository.UpdateRotated(ctx, UpdateRotatedParams{
		ID:        sessionID,
		Hash:      newHash,
		IP:        ip,
		ExpiresAt: tokens.RefreshExpiresAt,
	})
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Service) DeleteExpired(ctx context.Context) error {
	return s.repository.DeleteExpired(ctx)
}

func (s *Service) generateDeviceHash(deviceID, userID uuid.UUID) string {
	h := hmac.New(sha256.New, []byte(s.config.Auth.DeviceSalt))
	h.Write([]byte(deviceID.String() + ":" + userID.String()))

	return hex.EncodeToString(h.Sum(nil))
}
