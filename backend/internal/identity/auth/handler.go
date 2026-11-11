package auth

import (
	"context"
	"net/netip"
	"time"
	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/jethroTLOU/tokigura/internal/core/apperr"
	"github.com/jethroTLOU/tokigura/internal/core/config"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"github.com/jethroTLOU/tokigura/internal/core/validation"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
	"github.com/jethroTLOU/tokigura/internal/server/middleware"
	"go.uber.org/fx"
)

type service interface {
	Register(ctx context.Context, dto RegisterRequest) error
	Login(ctx context.Context, dto LoginDTO) (*jwt.TokenPair, error)
	Me(ctx context.Context, id uuid.UUID) (*user.User, error)
	Refresh(ctx context.Context, refreshToken string, ip netip.Addr) (*jwt.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type Handler struct {
	validator *validation.Validator
	service   service
	config    *config.Config
	authMW    *middleware.AuthMiddleware
}

type HandlerParams struct {
	fx.In

	Validator *validation.Validator
	Service   service
	Config    *config.Config
	AuthMW    *middleware.AuthMiddleware
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		validator: params.Validator,
		service:   params.Service,
		config:    params.Config,
		authMW:    params.AuthMW,
	}
}

func (h *Handler) Register(r fiber.Router) {
	group := r.Group("/auth")

	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(ctx fiber.Ctx) string {
			return ctx.IP()
		},
		LimitReached: func(ctx fiber.Ctx) error {
			return apperr.New(
				apperr.CodeBadRequest,
				"Too many requests",
			)
		},
	})

	group.Post("/register", authLimiter, h.register)
	group.Post("/login", authLimiter, h.login)

	group.Get("/me", h.authMW.Required(), h.getMe)
	group.Post("/refresh", h.refresh)
	group.Post("/logout", h.logout)
}

func (h *Handler) register(ctx fiber.Ctx) error {
	req := RegisterRequest{}

	if err := ctx.Bind().Body(&req); err != nil {
		return apperr.NewBindError(err)
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	err := h.service.Register(ctx.Context(), req)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) login(ctx fiber.Ctx) error {
	req := LoginRequest{}

	if err := ctx.Bind().Body(&req); err != nil {
		return apperr.NewBindError(err)
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	ip, err := netip.ParseAddr(ctx.IP())
	if err != nil {
		return apperr.New(apperr.CodeBadRequest, "invalid IP address")
	}

	dto := LoginDTO{
		Username:  req.Username,
		Password:  req.Password,
		IP:        ip,
		UserAgent: string(ctx.Request().Header.UserAgent()),
	}

	tokens, err := h.service.Login(ctx.Context(), dto)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Secure:   h.config.Auth.SecureCookies,
		Expires:  tokens.RefreshExpiresAt,
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/api/auth",
	})

	return ctx.JSON(&fiber.Map{
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) getMe(ctx fiber.Ctx) error {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return apperr.New(apperr.CodeInternal, "failed to get user id from context")
	}
	u, err := h.service.Me(ctx.Context(), userID)
	if err != nil {
		return err
	}

	return ctx.JSON(&MeResponse{
		Username: u.Username,
		ID:       u.ID,
		Email:    u.Email,
		Role:     u.Role,
		Avatar:   u.Avatar,
	})
}

func (h *Handler) refresh(ctx fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return apperr.New(apperr.CodeUnauthorized, "missing refresh token")
	}

	ip, err := netip.ParseAddr(ctx.IP())
	if err != nil {
		return apperr.New(apperr.CodeBadRequest, "invalid IP address")
	}

	tokens, err := h.service.Refresh(ctx.Context(), refreshToken, ip)
	if err != nil {
		ctx.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HTTPOnly: true,
			Secure:   h.config.Auth.SecureCookies,
			Expires:  time.Now().Add(-1 * time.Hour),
		})
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Secure:   h.config.Auth.SecureCookies,
		Expires:  tokens.RefreshExpiresAt,
	})

	return ctx.JSON(&fiber.Map{
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) logout(ctx fiber.Ctx) error {
	refreshToken := ctx.Cookies("refresh_token")

	if refreshToken != "" {
		_ = h.service.Logout(ctx.Context(), refreshToken)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   h.config.Auth.SecureCookies,
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
	})

	return ctx.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
