package auth

import (
	"context"
	"net/http"
	"net/netip"
	"time"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
	"github.com/jethroTLOU/tokigura/internal/core/config"
	"github.com/jethroTLOU/tokigura/internal/core/jwt"
	"github.com/jethroTLOU/tokigura/internal/core/validation"
	"github.com/jethroTLOU/tokigura/internal/identity/user"
	"github.com/jethroTLOU/tokigura/internal/server"
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

func (h *Handler) Register(mux *http.ServeMux) {
	authLimiter := middleware.RateLimit(5, time.Minute)

	mux.Handle("POST /auth/register", authLimiter(server.Wrap(h.register)))
	mux.Handle("POST /auth/login", authLimiter(server.Wrap(h.login)))
	mux.Handle("GET /auth/me", h.authMW.Required(server.Wrap(h.getMe)))
	mux.HandleFunc("POST /auth/refresh", server.Wrap(h.refresh))
	mux.HandleFunc("POST /auth/logout", server.Wrap(h.logout))
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) error {
	var req RegisterRequest
	if err := server.ReadJSON(w, r, &req); err != nil {
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	if err := h.service.Register(r.Context(), req); err != nil {
		return err
	}

	return server.WriteJSON(
		w,
		http.StatusCreated,
		server.Map{"message": "User created successfully"},
	)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := server.ReadJSON(w, r, &req); err != nil {
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	ip, err := middleware.ExtractIP(r)
	if err != nil {
		return err
	}

	dto := LoginDTO{
		Username:  req.Username,
		Password:  req.Password,
		IP:        ip,
		UserAgent: r.UserAgent(),
	}

	tokens, err := h.service.Login(r.Context(), dto)
	if err != nil {
		return err
	}

	h.setRefreshCookie(w, tokens.RefreshToken, tokens.RefreshExpiresAt)
	return server.WriteJSON(w, http.StatusOK, server.Map{
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) error {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		return apperr.New(apperr.CodeInternal, "failed to get user id from context")
	}

	u, err := h.service.Me(r.Context(), userID)
	if err != nil {
		return err
	}

	return server.WriteJSON(w, http.StatusOK, &MeResponse{
		Username: u.Username,
		ID:       u.ID,
		Email:    u.Email,
		Role:     u.Role,
		Avatar:   u.Avatar,
	})
}

func (h *Handler) refresh(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		return apperr.New(apperr.CodeUnauthorized, "missing refresh token")
	}

	ip, err := middleware.ExtractIP(r)
	if err != nil {
		return err
	}

	tokens, err := h.service.Refresh(r.Context(), cookie.Value, ip)
	if err != nil {
		h.clearRefreshCookie(w)
		return err
	}

	h.setRefreshCookie(w, tokens.RefreshToken, tokens.RefreshExpiresAt)
	return server.WriteJSON(w, http.StatusOK, server.Map{
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) error {
	if cookie, err := r.Cookie("refresh_token"); err == nil && cookie.Value != "" {
		_ = h.service.Logout(r.Context(), cookie.Value)
	}

	h.clearRefreshCookie(w)
	return server.WriteJSON(w, http.StatusOK, server.Map{
		"message": "Successfully logged out",
	})
}

func (h *Handler) setRefreshCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{ //nolint:gosec //  Secure flag is configured dynamically
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   h.config.Auth.SecureCookies,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth",
		Expires:  expires,
	})
}

func (h *Handler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ //nolint:gosec //  Secure flag is configured dynamically
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   h.config.Auth.SecureCookies,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth",
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
	})
}
