package service

import (
	"context"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Fetch(ctx context.Context, name string) error {
	return nil
}
func (s *Service) Sync(ctx context.Context, name string) error {
	return nil
}

func (s *Service) List(ctx context.Context, name string, minScore float64) (interface{}, error) {
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, name string) error {
	return nil
}

func (s *Service) Patch(ctx context.Context, name string) error {
	return nil
}
