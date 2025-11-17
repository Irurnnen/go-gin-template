package cat

import (
	"context"
)

type (
	CatRepositoryInterface interface {
		CatCreate(ctx context.Context, data *CatCreateRepo) error
		CatSearch(ctx context.Context, params *CatSearchParams) ([]*Cat, error)
		CatByID(ctx context.Context, id *CatID) (*Cat, error)
		CatExists(ctx context.Context, id *CatID) error
		CatChangeStatus(ctx context.Context, id *CatID, status *CatStatus) error
		CatDelete(ctx context.Context, id *CatID) error
		ToyCreate(ctx context.Context, catID *CatID, data *ToyCreateRepo) error
		ToyExists(ctx context.Context, id *ToyID) error
		ToyBelongs(ctx context.Context, catID *CatID, toyID *ToyID) error
		ToySearch(ctx context.Context, catID *CatID, params *ToySearchParams) ([]*Toy, error)
		ToyByID(ctx context.Context, toyID *ToyID) (*Toy, error)
		ToyDelete(ctx context.Context, toyID *ToyID) error
	}
)
