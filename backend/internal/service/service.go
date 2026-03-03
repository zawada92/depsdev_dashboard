package service

import (
	"context"

	"dependency-dashboard/internal/depsdev"
	"dependency-dashboard/internal/model"
	"dependency-dashboard/internal/repository"
)

type Service struct {
	repo   *repository.Repository
	client *depsdev.Client
}

func New(repo *repository.Repository, client *depsdev.Client) *Service {
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

func (s *Service) Patch(ctx context.Context, name string) error {
	// TODO_TOM
	return nil
}
