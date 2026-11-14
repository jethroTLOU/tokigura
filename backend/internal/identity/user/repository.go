package user

import (
	"context"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/database"
	"go.uber.org/fx"
)

type Repository struct {
	db *database.DB
}

type RepositoryParams struct {
	fx.In
	DB *database.DB
}

func NewRepository(repo RepositoryParams) *Repository {
	return &Repository{db: repo.DB}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	_, err := r.db.Conn(ctx).NewInsert().Model(u).Exec(ctx)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	u := new(User)
	err := r.db.Conn(ctx).NewSelect().Model(u).Where("username = ?", username).Scan(ctx)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := new(User)
	err := r.db.Conn(ctx).NewSelect().Model(u).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}
