package cat

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type (
	CatService struct {
		repo   CatRepositoryInterface
		logger *zerolog.Logger
	}
)

func NewCatService(repo CatRepositoryInterface, logger *zerolog.Logger) *CatService {
	return &CatService{
		repo:   repo,
		logger: logger,
	}
}

func (cs *CatService) CatCreate(ctx context.Context, data *CatCreate) (*CatID, error) {
	id := uuid.NewString()

	repoData := &CatCreateRepo{
		ID:             id,
		Name:           data.Name,
		Status:         CatStatusAvailable,
		Breed:          data.Breed,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
		BirthTimestamp: data.BirthTimestamp,
		Color:          data.Color,
		Gender:         data.Gender,
		Weight:         data.Weight,
		Description:    data.Description,
	}

	err := cs.repo.CatCreate(ctx, repoData)
	if err != nil {
		return nil, err
	}

	return &CatID{ID: id}, nil
}

func (cs *CatService) CatSearch(ctx context.Context, params *CatSearchParams) ([]*Cat, error) {
	cats, err := cs.repo.CatSearch(ctx, params)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (cs *CatService) CatByID(ctx context.Context, id *CatID) (*Cat, error) {
	cat, err := cs.repo.CatByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (cs *CatService) CatChangeStatus(ctx context.Context, id *CatID, status *CatStatusStruct) error {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, id); err != nil {
		return err
	}

	// Change status
	err := cs.repo.CatChangeStatus(ctx, id, &status.Status)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CatService) CatDelete(ctx context.Context, id *CatID) error {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, id); err != nil {
		return err
	}

	// Change status
	err := cs.repo.CatDelete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CatService) ToyCreate(ctx context.Context, catID *CatID, data *ToyCreate) (*ToyID, error) {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, catID); err != nil {
		return nil, err
	}

	// Create id
	id := uuid.NewString()

	// Create structure for create ToyID
	repoData := &ToyCreateRepo{
		ID:        id,
		Name:      data.Name,
		Type:      data.Type,
		Color:     data.Color,
		Material:  data.Material,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Create Toy
	err := cs.repo.ToyCreate(ctx, catID, repoData)
	if err != nil {
		return nil, err
	}

	return &ToyID{ID: id}, nil
}

func (cs *CatService) ToySearch(ctx context.Context, catID *CatID, params *ToySearchParams) ([]*Toy, error) {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, catID); err != nil {
		return nil, err
	}

	// Get toys by cat's id
	toys, err := cs.repo.ToySearch(ctx, catID, params)
	if err != nil {
		return nil, err
	}
	return toys, nil
}

func (cs *CatService) ToyByID(ctx context.Context, catID *CatID, toyID *ToyID) (*Toy, error) {
	// Check cat exists
	if err := cs.repo.CatExists(ctx, catID); err != nil {
		return nil, err
	}

	// Check toy exists
	if err := cs.repo.ToyExists(ctx, toyID); err != nil {
		return nil, err
	}

	// Check toy belongs to cat
	if err := cs.repo.ToyBelongs(ctx, catID, toyID); err != nil {
		return nil, err
	}

	// Get toy by toy id
	toy, err := cs.repo.ToyByID(ctx, toyID)
	if err != nil {
		return nil, err
	}

	return toy, nil
}

func (cs *CatService) ToyDelete(ctx context.Context, catID *CatID, toyID *ToyID) error {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, catID); err != nil {
		return err
	}

	// Check toy exists
	if err := cs.repo.ToyExists(ctx, toyID); err != nil {
		return err
	}

	// Check toy belongs to cat
	if err := cs.repo.ToyBelongs(ctx, catID, toyID); err != nil {
		return err
	}

	// Delete cat
	err := cs.repo.ToyDelete(ctx, toyID)
	if err != nil {
		return err
	}

	return nil
}
