package services

import (
	"github.com/Irurnnen/go-gin-template/internal/repository"
	"go.uber.org/zap"
)

type HelloService struct {
	repo   repository.HelloRepositoryInterface
	logger *zap.Logger
}

type HelloServiceInterface interface {
	GetHelloMessage() (string, error)
}

func NewHelloService(repo repository.HelloRepositoryInterface, logger *zap.Logger) *HelloService {
	return &HelloService{
		repo:   repo,
		logger: logger,
	}
}

func (s *HelloService) GetHelloMessage() (string, error) {
	message, err := s.repo.GetHelloMessage()
	if err != nil {
		return "", err
	}
	return message, nil
}
