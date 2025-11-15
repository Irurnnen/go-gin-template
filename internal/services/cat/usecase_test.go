package cat

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCatService(t *testing.T) {
	// Repo and logger mocks
	logger := zerolog.Nop()
	mockRepo := NewMockCatRepositoryInterface(t)

	// Create service
	service := NewCatService(mockRepo, &logger)

	// Assertions
	if assert.NotNil(t, service) {
		assert.NotEmpty(t, service)
	}
}

func TestCatService_Create_Success(t *testing.T) {
	ctx := context.Background()

	// input data
	input := &CreateCat{
		Name:           "Tom",
		Breed:          "Tabby",
		Color:          "gray",
		Gender:         "male",
		Weight:         4,
		Description:    "friendly",
		BirthTimestamp: time.Now(),
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	var captured *CreateCatRepo
	mockRepo.EXPECT().Create(ctx, mock.Anything).
		Run(func(_ctx context.Context, data *CreateCatRepo) {
			captured = data
		}).Return(nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	id, err := service.Create(ctx, input)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, id) {
		assert.NotEmpty(t, id.ID)
	}

	// Verify repository received expected transformed data
	if assert.NotNil(t, captured) {
		assert.NotEmpty(t, captured.ID)
		assert.Equal(t, input.Name, captured.Name)
		assert.Equal(t, input.Breed, captured.Breed)
		assert.Equal(t, input.Color, captured.Color)
		assert.Equal(t, input.Gender, captured.Gender)
		assert.Equal(t, input.Weight, captured.Weight)
		assert.Equal(t, input.Description, captured.Description)
		assert.Equal(t, input.BirthTimestamp, captured.BirthTimestamp)
		assert.Equal(t, CatStatusAvailable, captured.Status)
		assert.False(t, captured.CreatedAt.IsZero())
		assert.False(t, captured.UpdatedAt.IsZero())
	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_Create_RepoError(t *testing.T) {
	ctx := context.Background()

	// Input data
	input := &CreateCat{
		Name:           "Jerry",
		Breed:          "Siamese",
		BirthTimestamp: time.Now(),
		Color:          "green",
		Gender:         "male",
		Weight:         0,
		Description:    "asdasd",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		Create(ctx, mock.Anything).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	id, err := service.Create(ctx, input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, id)

	mockRepo.AssertExpectations(t)
}

func TestCatService_Search_Success(t *testing.T) {
	ctx := context.Background()

	// input data
	name := "Tom"
	breed := "Tabby"
	params := &CatSearchParams{
		Name:  &name,
		Breed: &breed,
	}

	// output data
	cats := []*Cat{
		{Name: "Hello"},
		{Name: "World"},
		{Name: "One"},
		{Name: "Two"},
		{Name: "Three"},
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	mockRepo.EXPECT().Search(ctx, params).
		Return(cats, nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cats, err := service.Search(ctx, params)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, cats) {
		for i, cat := range cats {
			assert.NotEmpty(t, cat.Name, "i=%d got incorrect cat", i)
		}
	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_Search_RepoErr(t *testing.T) {
	ctx := context.Background()

	// input data
	name := "Tom"
	breed := "Tabby"
	params := &CatSearchParams{
		Name:  &name,
		Breed: &breed,
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	mockRepo.EXPECT().
		Search(ctx, params).
		Return(nil, errors.New("db error")).Once()

	// Cat Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cats, err := service.Search(ctx, params)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, cats)

	mockRepo.AssertExpectations(t)
}

func TestCatService_GetByID_Success(t *testing.T) {
	ctx := context.Background()

	// input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	exceptedCat := &Cat{
		Name:        "Tom",
		Breed:       "Tabby",
		Color:       "gray",
		Gender:      "male",
		Weight:      4,
		Description: "friendly",
		ID:          "00000000-0000-0000-0000-000000000000",
		Status:      CatStatusAvailable,
		Age:         12,
		CreatedAt:   int(time.Now().Unix()) - 1000,
		UpdatedAt:   int(time.Now().Unix()),
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	mockRepo.EXPECT().GetByID(ctx, id).
		Return(exceptedCat, nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cat, err := service.GetByID(ctx, id)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, cat) {
		assert.Equal(t, exceptedCat.ID, cat.ID)
		assert.Equal(t, exceptedCat.Name, cat.Name)
		assert.Equal(t, exceptedCat.Status, cat.Status)
		assert.Equal(t, exceptedCat.Color, cat.Color)
		assert.Equal(t, exceptedCat.Gender, cat.Gender)
		assert.Equal(t, exceptedCat.Weight, cat.Weight)
		assert.Equal(t, exceptedCat.Description, cat.Description)
		assert.Equal(t, exceptedCat.Breed, cat.Breed)
		assert.Equal(t, exceptedCat.Age, cat.Age)
		assert.Equal(t, exceptedCat.CreatedAt, cat.CreatedAt)
		assert.Equal(t, exceptedCat.UpdatedAt, cat.UpdatedAt)
	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_GetByID_RepoErr(t *testing.T) {
	ctx := context.Background()

	// input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		GetByID(ctx, id).
		Return(nil, errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cat, err := service.GetByID(ctx, id)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, cat)

	mockRepo.AssertExpectations(t)
}
func TestCatService_ChangeStatus_Success(t *testing.T) {
	ctx := context.Background()

	// input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	status := CatStatusSick

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	mockRepo.EXPECT().CatExists(ctx, id).
		Return(nil).Once()
	mockRepo.EXPECT().ChangeStatus(ctx, id, &status).
		Return(nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.ChangeStatus(ctx, id, &status)

	// Assertions
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_ChangeStatus_CatExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	status := CatStatusSick

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, id).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.ChangeStatus(ctx, id, &status)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
func TestCatService_ChangeStatus_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	status := CatStatusSick

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, id).
		Return(nil).Once()
	mockRepo.EXPECT().
		ChangeStatus(ctx, id, &status).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.ChangeStatus(ctx, id, &status)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_Delete_Success(t *testing.T) {
	ctx := context.Background()

	// input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	mockRepo.EXPECT().CatExists(ctx, id).
		Return(nil).Once()
	mockRepo.EXPECT().Delete(ctx, id).
		Return(nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.Delete(ctx, id)

	// Assertions
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_Delete_CatExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, id).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.Delete(ctx, id)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_Delete_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, id).
		Return(nil).Once()
	mockRepo.EXPECT().
		Delete(ctx, id).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.Delete(ctx, id)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_AddToy_Success(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}
	input := &CreateToy{
		Name:     "Gilbert",
		Type:     "bear",
		Color:    "white",
		Material: "synthetic",
	}

	// Mock repository
	var captured *CreateToyRepo
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		AddToy(ctx, catID, mock.Anything).
		Run(func(ctx context.Context, catID *CatID, data *CreateToyRepo) {
			captured = data
		}).
		Return(nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toyID, err := service.AddToy(ctx, catID, input)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, toyID) {
		assert.NotEmpty(t, toyID.ID)
	}

	// Verify repository received expected transformed data
	if assert.NotNil(t, captured) {
		assert.NotEmpty(t, captured.ID)
		assert.Equal(t, input.Name, captured.Name)
		assert.Equal(t, input.Type, captured.Type)
		assert.Equal(t, input.Color, captured.Color)
		assert.Equal(t, input.Material, captured.Material)
		assert.False(t, captured.CreatedAt.IsZero())
		assert.False(t, captured.UpdatedAt.IsZero())
	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_AddToy_CatExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	input := &CreateToy{
		Name:     "Gilbert",
		Type:     "bear",
		Color:    "white",
		Material: "synthetic",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, id).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	_, err := service.AddToy(ctx, id, input)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
func TestCatService_AddToy_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}
	input := &CreateToy{
		Name:     "Gilbert",
		Type:     "bear",
		Color:    "white",
		Material: "synthetic",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		AddToy(ctx, catID, mock.Anything).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toyID, err := service.AddToy(ctx, catID, input)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, toyID)

	mockRepo.AssertExpectations(t)
}

func TestCatService_GetToys_Success(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	exceptedToys := []*Toy{
		{Name: "Hello"},
		{Name: "World"},
		{Name: "One"},
		{Name: "Two"},
		{Name: "Three"},
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		GetToys(ctx, catID).
		Return(exceptedToys, nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toys, err := service.GetToys(ctx, catID)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, toys) {
		for i, toy := range toys {
			assert.Equal(t, toy.Name, exceptedToys[i].Name,
				"i=%d got incorrect toy", i)
		}
	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_GetToys_CatExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, id).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	_, err := service.GetToys(ctx, id)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_GetToys_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		GetToys(ctx, catID).
		Return(nil, errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toyID, err := service.GetToys(ctx, catID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, toyID)

	mockRepo.AssertExpectations(t)
}

func TestCatService_GetToyByID_Success(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	exceptedToy := &Toy{
		Name:      "Hello",
		ID:        "00000000-0000-0000-0000-000000000000",
		Type:      "bear",
		Color:     "green",
		Material:  "synthetic",
		CreatedAt: int(time.Now().Unix()),
		UpdatedAt: int(time.Now().Unix()) - 36000,
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyBelongs(ctx, catID, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		GetToyByID(ctx, toyID).
		Return(exceptedToy, nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toy, err := service.GetToyByID(ctx, catID, toyID)

	// Assertions
	assert.NoError(t, err)
	if assert.NotNil(t, toy) {
		assert.Equal(t, toy.Name, exceptedToy.Name)
		assert.Equal(t, toy.ID, exceptedToy.ID)
		assert.Equal(t, toy.CreatedAt, exceptedToy.CreatedAt)
		assert.Equal(t, toy.Type, exceptedToy.Type)
		assert.Equal(t, toy.Color, exceptedToy.Color)
		assert.Equal(t, toy.Material, exceptedToy.Material)
		assert.Equal(t, toy.UpdatedAt, exceptedToy.UpdatedAt)

	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_GetToyByID_CatExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toy, err := service.GetToyByID(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, toy)

	mockRepo.AssertExpectations(t)
}
func TestCatService_GetToyByID_ToyExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toy, err := service.GetToyByID(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, toy)

	mockRepo.AssertExpectations(t)
}
func TestCatService_GetToyByID_ToyBelongs_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyBelongs(ctx, catID, toyID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toy, err := service.GetToyByID(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, toy)

	mockRepo.AssertExpectations(t)
}
func TestCatService_GetToyByID_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyBelongs(ctx, catID, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		GetToyByID(ctx, toyID).
		Return(nil, errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toy, err := service.GetToyByID(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, toy)

	mockRepo.AssertExpectations(t)
}

func TestCatService_DeleteToy_Success(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyBelongs(ctx, catID, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		DeleteToy(ctx, toyID).
		Return(nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.DeleteToy(ctx, catID, toyID)

	// Assertions
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_DeleteToy_CatExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.DeleteToy(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
func TestCatService_DeleteToy_ToyExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.DeleteToy(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_DeleteToy_ToyBelongs_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyBelongs(ctx, catID, toyID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.DeleteToy(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCatService_DeleteToy_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	catID := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	toyID := &ToyID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyExists(ctx, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyBelongs(ctx, catID, toyID).
		Return(nil).Once()
	mockRepo.EXPECT().
		DeleteToy(ctx, toyID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.DeleteToy(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
