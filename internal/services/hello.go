package services

import (
	"context"

	"github.com/Irurnnen/go-gin-template/internal/repository"
	"github.com/rs/zerolog"
)

type (
	HelloService struct {
		repo   repository.HelloRepositoryInterface
		logger *zerolog.Logger
	}

	HelloServiceInterface interface {
		GetHelloMessage(context.Context) (string, error)
	}
)

func NewHelloService(repo repository.HelloRepositoryInterface, logger *zerolog.Logger) *HelloService {
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
