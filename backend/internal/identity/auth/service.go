package auth

import (
	"context"
	"net/netip"
	"strings"
	"time"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
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
	RotateSession(
		ctx context.Context,
		sessionID, userID uuid.UUID,
		role user.Role,
		ip netip.Addr,
	) (*jwt.TokenPair, error)
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
		ID:       uuid.New(),
		Email:    strings.ToLower(dto.Email),
		Username: strings.ToLower(dto.Username),
		Password: string(hashedPassword),
		Role:     user.UserRole,
	}

	if err := s.users.Create(ctx, newUser); err != nil {
		return apperr.Wrap(err, apperr.CodeConflict, "username or email already exists")
	}

	return nil
}

func (s *Service) Login(ctx context.Context, dto LoginDTO) (*jwt.TokenPair, error) {
	u, err := s.users.FindByUsername(ctx, dto.Username)
	if err != nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(dto.Password)); err != nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid username or password")
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
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, apperr.New(apperr.CodeNotFound, "user not found")
	}

	return u, nil
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
	if sess == nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "invalid session")
	}

	if sess.ExpiresAt.Before(time.Now()) {
		_ = s.sessions.DeleteByToken(ctx, refreshToken)
		return nil, apperr.New(apperr.CodeUnauthorized, "session expired")
	}

	u, err := s.users.FindByID(ctx, sess.UserID)
	if err != nil {
		return nil, apperr.New(apperr.CodeUnauthorized, "user not found")
	}

	return s.sessions.RotateSession(ctx, sess.ID, u.ID, u.Role, ip)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.sessions.DeleteByToken(ctx, refreshToken)
}
