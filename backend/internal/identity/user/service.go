package user

import (
	"context"
	"uuid"
)

type repository interface {
	Create(ctx context.Context, u *User) error
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, u *User) error {
	return s.repo.Create(ctx, u)
}

func (s *Service) FindByUsername(ctx context.Context, username string) (*User, error) {
	return s.repo.FindByUsername(ctx, username)
}

func (s *Service) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.repo.FindByID(ctx, id)
}
