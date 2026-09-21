package user

import (
	"context"
	_ "embed"
	"uuid"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Repository struct {
	db *pgxpool.Pool
}

type RepositoryParams struct {
	fx.In

	DB *pgxpool.Pool
}

func NewRepository(repo RepositoryParams) *Repository {
	return &Repository{db: repo.DB}
}

var (
	//go:embed sql/create.sql
	queryCreate string

	//go:embed sql/find-by-username.sql
	queryFindByUsername string

	//go:embed sql/find-by-id.sql
	queryFindByID string
)

func (r *Repository) Create(ctx context.Context, u *User) error {
	_, err := r.db.Exec(ctx, queryCreate,
		u.Email, u.Username, u.Role, u.Password, u.Avatar,
	)
	return err
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	u := new(User)

	err := pgxscan.Get(ctx, r.db, u, queryFindByUsername, username)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u := new(User)

	err := pgxscan.Get(ctx, r.db, u, queryFindByID, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}
