package hello

import (
	"context"

	"github.com/rs/zerolog"
)

type (
	HelloService struct {
		repo   HelloRepositoryInterface
		logger *zerolog.Logger
	}
)

func NewHelloService(repo HelloRepositoryInterface, logger *zerolog.Logger) *HelloService {
	return &HelloService{
		repo:   repo,
		logger: logger,
	}
}

func (s *HelloService) GetHelloMessage(ctx context.Context) (*Message, error) {
	message, err := s.repo.GetHelloMessage(ctx)
	if err != nil {
		return nil, err
	}
	return message, nil
}
