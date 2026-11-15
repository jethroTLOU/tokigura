package auth

import (
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

func (r *RegisterRequest) Validate() error {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Username = strings.ToLower(strings.TrimSpace(r.Username))

	errs := make(map[string]any)

	if r.Email == "" {
		errs["email"] = "email is required"
	} else if !isValidEmail(r.Email) {
		errs["email"] = "invalid email format"
	}

	if msg := validateUsername(r.Username); msg != "" {
		errs["username"] = msg
	}

	if msg := validatePassword(r.Password); msg != "" {
		errs["password"] = msg
	}

	if len(errs) > 0 {
		return apperr.NewValidation(errs)
	}
	return nil
}

func (r *LoginRequest) Validate() error {
	r.Username = strings.ToLower(strings.TrimSpace(r.Username))

	errs := make(map[string]any)

	if msg := validateUsername(r.Username); msg != "" {
		errs["username"] = msg
	}

	if msg := validatePassword(r.Password); msg != "" {
		errs["password"] = msg
	}

	if len(errs) > 0 {
		return apperr.NewValidation(errs)
	}
	return nil
}

func validateUsername(username string) string {
	if username == "" {
		return "username is required"
	}

	length := utf8.RuneCountInString(username)
	if length < 3 || length > 32 {
		return "username must be between 3 and 32 characters"
	}

	for _, ch := range username {
		if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') {
			return "username must contain only letters and numbers"
		}
	}
	return ""
}

func validatePassword(password string) string {
	if password == "" {
		return "password is required"
	}

	if utf8.RuneCountInString(password) < 8 {
		return "password must be at least 8 characters"
	}

	if len(password) > 72 {
		return "password must not exceed 72 bytes"
	}

	return ""
}

func isValidEmail(email string) bool {
	if len(email) > 254 {
		return false
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return false
	}

	_, domain, ok := strings.Cut(email, "@")
	return ok && strings.Contains(domain, ".") && !strings.HasPrefix(domain, ".") &&
		!strings.HasSuffix(domain, ".")
}
