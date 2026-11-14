package auth

import (
	"context"
	"errors"
	"net/netip"
	"strings"
	"time"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
	"github.com/jethroTLOU/tokigura/internal/core/domain"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"github.com/jethroTLOU/tokigura/internal/identity/session"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
	"golang.org/x/crypto/bcrypt"
)

type userService interface {
	Create(ctx context.Context, u *user.User) error
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*user.User, error)
}

type sessions interface {
	Create(ctx context.Context, dto session.CreateSessionDTO) (*jwt.TokenPair, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*session.Session, error)
	DeleteByToken(ctx context.Context, refreshToken string) error
	RotateSession(ctx context.Context, dto session.RotateSessionDTO) (*jwt.TokenPair, error)
}

type Service struct {
	users    userService
	sessions sessions
}

func NewService(users userService, sessions sessions) *Service {
	return &Service{
		users:    users,
		sessions: sessions,
	}
}

func (s *Service) Register(ctx context.Context, dto RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), 12)
	if err != nil {
		return apperr.Wrap(err, apperr.CodeInternal, "failed to hash password")
	}

	newUser := &user.User{
		ID:       uuid.NewV7(),
		Email:    strings.ToLower(dto.Email),
		Username: strings.ToLower(dto.Username),
		Password: string(hashedPassword),
		Role:     domain.RoleUser,
	}

	return s.users.Create(ctx, newUser)
}

func (s *Service) Login(ctx context.Context, dto LoginDTO) (*jwt.TokenPair, error) {
	u, err := s.users.FindByUsername(ctx, strings.ToLower(dto.Username))
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	sessionDTO := session.CreateSessionDTO{
		UserID:    u.ID,
		Role:      u.Role,
		IP:        dto.IP,
		UserAgent: dto.UserAgent,
	}

	return s.sessions.Create(ctx, sessionDTO)
}

func (s *Service) Me(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return s.users.FindByID(ctx, userID)
}

func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
	ip netip.Addr,
) (*jwt.TokenPair, error) {
	sess, err := s.sessions.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if sess.ExpiresAt.Before(time.Now()) {
		_ = s.sessions.DeleteByToken(ctx, refreshToken)
		return nil, ErrSessionExpired
	}

	u, err := s.users.FindByID(ctx, sess.UserID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return s.sessions.RotateSession(ctx, session.RotateSessionDTO{
		SessionID: sess.ID,
		UserID:    u.ID,
		Role:      u.Role,
		IP:        ip,
		OldHash:   sess.TokenHash,
	})
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.sessions.DeleteByToken(ctx, refreshToken)
}
