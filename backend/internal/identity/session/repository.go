package session

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"time"
	"uuid"

	"github.com/jethroTLOU/tokigura/internal/core/apperr"
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

func (r *Repository) Create(ctx context.Context, s *Session) error {
	_, err := r.db.Conn(ctx).NewInsert().Model(s).Exec(ctx)
	return err
}

func (r *Repository) GetByTokenHash(ctx context.Context, hash string) (*Session, error) {
	s := new(Session)

	if err := r.db.Conn(ctx).
		NewSelect().
		Model(s).
		Where("token_hash = ?", hash).
		Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return s, nil
}

func (r *Repository) DeleteByTokenHash(ctx context.Context, hash string) error {
	_, err := r.db.Conn(ctx).
		NewDelete().
		Model((*Session)(nil)).
		Where("token_hash = ?", hash).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.
		Conn(ctx).
		NewDelete().
		Model((*Session)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)

	return err
}

func (r *Repository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.Conn(ctx).
		NewDelete().
		Model((*Session)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)

	return err
}

func (r *Repository) UpdateRotated(ctx context.Context, params UpdateRotatedParams) error {
	res, err := r.db.Conn(ctx).
		NewUpdate().
		Model((*Session)(nil)).
		Set("token_hash = ?", params.NewHash).
		Set("last_ip_address = ?", params.IP).
		Set("expires_at = ?", params.ExpiresAt).
		Where("id = ? AND token_hash = ?", params.ID, params.OldHash).
		Exec(ctx)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperr.New(apperr.CodeUnauthorized, "session not found or expired")
	}

	return nil
}
