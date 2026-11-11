package session

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"
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

func (r *Repository) Create(ctx context.Context, s *Session) error {
	_, err := r.db.NewInsert().Model(s).Exec(ctx)
	return err
}

func (r *Repository) GetByTokenHash(ctx context.Context, hash string) (*Session, error) {
	s := new(Session)

	if err := r.db.NewSelect().Model(s).Where("token_hash = ?", hash).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return s, nil
}

func (r *Repository) DeleteByTokenHash(ctx context.Context, hash string) error {
	_, err := r.db.NewDelete().
		Model((*Session)(nil)).
		Where("token_hash = ?", hash).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.NewDelete().
		Model((*Session)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*Session)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)

	return err
}

func (r *Repository) UpdateRotated(ctx context.Context, params UpdateRotatedParams) error {
	res, err := r.db.NewUpdate().
		Model((*Session)(nil)).
		Set("token_hash = ?", params.Hash).
		Set("last_ip_address = ?", params.IP).
		Set("expires_at = ?", params.ExpiresAt).
		Where("id = ?", params.ID).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return nil
	}

	return nil
}
