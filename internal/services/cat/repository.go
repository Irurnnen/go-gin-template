package cat

import (
	"context"
)

type (
	CatRepositoryInterface interface {
		Create(ctx context.Context, data *CreateCatRepo) error
		Search(ctx context.Context, params *CatSearchParams) ([]*Cat, error)
		GetByID(ctx context.Context, id *CatID) (*Cat, error)
		CatExists(ctx context.Context, id *CatID) error
		ChangeStatus(ctx context.Context, id *CatID, status *CatStatus) error
		Delete(ctx context.Context, id *CatID) error
		AddToy(ctx context.Context, catID *CatID, data *CreateToyRepo) error
		ToyExists(ctx context.Context, id *ToyID) error
		ToyBelongs(ctx context.Context, catID *CatID, toyID *ToyID) error
		GetToys(ctx context.Context, catID *CatID) ([]*Toy, error)
		GetToyByID(ctx context.Context, toyID *ToyID) (*Toy, error)
		DeleteToy(ctx context.Context, toyID *ToyID) error
	}
)
