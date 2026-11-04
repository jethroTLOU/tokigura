package user

import (
	"context"
	_ "embed"
	"uuid"

	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type Repository struct {
	db *bun.DB
}

type RepositoryParams struct {
	fx.In

	DB *bun.DB
}

func NewRepository(repo RepositoryParams) *Repository {
	return &Repository{db: repo.DB}
}

func (r *Repository) Create(ctx context.Context, u *User) error {
	_, err := r.db.NewInsert().Model(u).Exec(ctx)
	return err
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	u := new(User)

	if err := r.db.NewSelect().Model(u).Where("username = ?", username).Scan(ctx); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := new(User)

	if err := r.db.NewSelect().Model(u).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return u, nil
}
