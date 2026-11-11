package user

import (
	"time"
	"uuid"

	"github.com/uptrace/bun"
)

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID `bun:"id,pk,type:text"`
	Email     string    `bun:"email,notnull,unique"`
	Username  string    `bun:"username,notnull,unique"`
	Role      Role      `bun:"role,notnull,type:text"`
	Password  string    `bun:"password,notnull"`
	Avatar    *string   `bun:"avatar"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
