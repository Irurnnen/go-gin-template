package services

import (
	"context"

	"github.com/rs/zerolog"
)

type (
	HelloService struct {
		repo   HelloRepositoryInterface
		logger *zerolog.Logger
	}

	HelloServiceInterface interface {
		GetHelloMessage(context.Context) (string, error)
	}

	HelloRepositoryInterface interface {
		GetHelloMessage(context.Context) (string, error)
	}
)

func NewHelloService(repo HelloRepositoryInterface, logger *zerolog.Logger) *HelloService {
	return &HelloService{
		repo:   repo,
		logger: logger,
	}
}

func (s *HelloService) GetHelloMessage(ctx context.Context) (string, error) {
	message, err := s.repo.GetHelloMessage(ctx)
	if err != nil {
		return "", err
	}
	return message, nil
}
