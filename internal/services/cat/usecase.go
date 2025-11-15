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

func (cs *CatService) Create(ctx context.Context, data *CreateCat) (*CatID, error) {
	id := uuid.NewString()

	repoData := &CreateCatRepo{
		ID:             id,
		Name:           data.Name,
		Status:         CatStatusAvailable,
		Breed:          data.Breed,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		BirthTimestamp: time.Time{},
		Color:          data.Color,
		Gender:         data.Gender,
		Weight:         data.Weight,
		Description:    data.Description,
	}

	err := cs.repo.Create(ctx, repoData)
	if err != nil {
		return nil, err
	}

	return &CatID{ID: id}, nil
}

func (cs *CatService) Search(ctx context.Context, params *CatSearchParams) ([]*Cat, error) {
	cats, err := cs.repo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (cs *CatService) GetByID(ctx context.Context, id *CatID) (*Cat, error) {
	cat, err := cs.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (cs *CatService) ChangeStatus(ctx context.Context, id *CatID, status *CatStatus) error {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, id); err != nil {
		return err
	}

	// Change status
	err := cs.repo.ChangeStatus(ctx, id, status)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CatService) Delete(ctx context.Context, id *CatID) error {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, id); err != nil {
		return err
	}

	// Change status
	err := cs.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CatService) AddToy(ctx context.Context, catID *CatID, data *CreateToy) (*ToyID, error) {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, catID); err != nil {
		return nil, err
	}

	// Create id
	id := uuid.NewString()

	// Create structure for create ToyID
	repoData := &CreateToyRepo{
		ID:        id,
		Name:      data.Name,
		Type:      data.Type,
		Color:     data.Color,
		Material:  data.Material,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create Toy
	err := cs.repo.AddToy(ctx, catID, repoData)
	if err != nil {
		return nil, err
	}

	return &ToyID{ID: id}, nil
}

func (cs *CatService) GetToys(ctx context.Context, catID *CatID) ([]*Toy, error) {
	// Check exists cat
	if err := cs.repo.CatExists(ctx, catID); err != nil {
		return nil, err
	}

	// Get toys by cat's id
	toys, err := cs.repo.GetToys(ctx, catID)
	if err != nil {
		return nil, err
	}
	return toys, nil
}

func (cs *CatService) GetToyByID(ctx context.Context, catID *CatID, toyID *ToyID) (*Toy, error) {
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
	toy, err := cs.repo.GetToyByID(ctx, toyID)
	if err != nil {
		return nil, err
	}

	return toy, nil
}

func (cs *CatService) DeleteToy(ctx context.Context, catID *CatID, toyID *ToyID) error {
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
	err := cs.repo.DeleteToy(ctx, toyID)
	if err != nil {
		return err
	}

	return nil
}
