package service

import (
	"context"
	"time"

	"dependency-dashboard/internal/domain"
	"dependency-dashboard/internal/model"
)

type Repository interface {
	Upsert(ctx context.Context, d *model.Dependency) error
	Update(ctx context.Context, d *model.Dependency) error
	List(ctx context.Context, name string, minScore float64) ([]model.Dependency, error)
	Delete(ctx context.Context, name string) error
	GetByName(ctx context.Context, name string) (*model.Dependency, error)
}

type DepsClient interface {
	Fetch(ctx context.Context, name string) (*model.Dependency, error)
}

type Service struct {
	repo   Repository
	client DepsClient
}

func New(repo Repository, client DepsClient) *Service {
	return &Service{repo: repo, client: client}
}

func (s *Service) Fetch(ctx context.Context, name string) error {
	dep, err := s.client.Fetch(ctx, name)
	if err != nil {
		return err
	}
	return s.repo.Upsert(ctx, dep)
}

func (s *Service) List(ctx context.Context, name string, minScore float64) ([]model.Dependency, error) {
	return s.repo.List(ctx, name, minScore)
}

func (s *Service) Delete(ctx context.Context, name string) error {
	return s.repo.Delete(ctx, name)
}

func (s *Service) Patch(
	ctx context.Context,
	name string,
	patch model.PatchDependencyRequest,
) (*model.Dependency, error) {

	existing, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if patch.Version != nil {
		if *patch.Version == "" {
			return nil, domain.ErrInvalidInput
		}
		existing.Version = *patch.Version
	}

	if patch.OpenSSFScore != nil {
		if *patch.OpenSSFScore < 0 || *patch.OpenSSFScore > 10 {
			return nil, domain.ErrInvalidInput
		}
		existing.OpenSSFScore = *patch.OpenSSFScore
	}

	existing.LastUpdated = time.Now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}
