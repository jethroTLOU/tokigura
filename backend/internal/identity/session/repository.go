package session

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

	//go:embed sql/get-by-token-hash.sql
	queryGetByTokenHash string

	//go:embed sql/delete-by-token-hash.sql
	queryDeleteByTokenHash string

	//go:embed sql/delete-all-by-user-id.sql
	queryDeleteAllByUserID string

	//go:embed sql/delete-expired.sql
	queryDeleteExpired string

	//go:embed sql/update-rotated.sql
	queryUpdateRotated string
)

func (r *Repository) Create(ctx context.Context, s *Session) error {
	_, err := r.db.Exec(ctx, queryCreate,
		s.UserID, s.TokenHash, s.DeviceInfo, s.DeviceHash,
		s.IPAddress, s.LastIPAddress, s.ExpiresAt,
	)
	return err
}

func (r *Repository) GetByTokenHash(ctx context.Context, hash string) (*Session, error) {
	s := new(Session)

	err := pgxscan.Get(ctx, r.db, s, queryGetByTokenHash, hash)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

func (r *Repository) DeleteByTokenHash(ctx context.Context, hash string) error {
	_, err := r.db.Exec(ctx, queryDeleteByTokenHash, hash)
	return err
}

func (r *Repository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, queryDeleteAllByUserID, userID)
	return err
}

func (r *Repository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.Exec(ctx, queryDeleteExpired)
	return err
}

func (r *Repository) UpdateRotated(ctx context.Context, params UpdateRotatedParams) error {
	_, err := r.db.Exec(
		ctx,
		queryUpdateRotated,
		params.Hash,
		params.IP,
		params.ExpiresAt,
		params.ID,
	)
	return err
}
