package hello

import "context"

type (
	HelloRepositoryInterface interface {
		GetHelloMessage(context.Context) (*Message, error)
	}
)
