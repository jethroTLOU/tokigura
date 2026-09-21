package user

import (
	"time"
	"uuid"
)

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Username  string    `db:"username"`
	Role      Role      `db:"role"`
	Password  string    `db:"password"`
	Avatar    *string   `db:"avatar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
