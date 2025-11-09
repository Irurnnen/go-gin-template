package hello

import "context"

type (
	HelloRepositoryInterface interface {
		GetHelloMessage(context.Context) (string, error)
	}
)
