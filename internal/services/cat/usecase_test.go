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
	description := "friendly"
	input := &CatCreate{
		Name:           "Tom",
		Breed:          "Tabby",
		Color:          "gray",
		Gender:         "male",
		Weight:         4,
		Description:    &description,
		BirthTimestamp: time.Now().Unix(),
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	var captured *CatCreateRepo
	mockRepo.EXPECT().CatCreate(ctx, mock.Anything).
		Run(func(_ctx context.Context, data *CatCreateRepo) {
			captured = data
		}).Return(nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	id, err := service.CatCreate(ctx, input)

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
		assert.Positive(t, captured.CreatedAt)
		assert.Positive(t, captured.UpdatedAt)
	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_Create_RepoError(t *testing.T) {
	ctx := context.Background()

	// Input data
	description := "asdasd"
	input := &CatCreate{
		Name:           "Jerry",
		Breed:          "Siamese",
		BirthTimestamp: time.Now().Unix(),
		Color:          "green",
		Gender:         "male",
		Weight:         0,
		Description:    &description,
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatCreate(ctx, mock.Anything).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	id, err := service.CatCreate(ctx, input)

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

	mockRepo.EXPECT().CatSearch(ctx, params).
		Return(cats, nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cats, err := service.CatSearch(ctx, params)

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
		CatSearch(ctx, params).
		Return(nil, errors.New("db error")).Once()

	// Cat Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cats, err := service.CatSearch(ctx, params)

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
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	mockRepo.EXPECT().CatByID(ctx, id).
		Return(exceptedCat, nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cat, err := service.CatByID(ctx, id)

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
		CatByID(ctx, id).
		Return(nil, errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	cat, err := service.CatByID(ctx, id)

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

	status := CatStatusStruct{
		Status: CatStatusSick,
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)

	mockRepo.EXPECT().CatExists(ctx, id).
		Return(nil).Once()
	mockRepo.EXPECT().CatChangeStatus(ctx, id, &status.Status).
		Return(nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.CatChangeStatus(ctx, id, &status)

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

	status := CatStatusStruct{
		Status: CatStatusSick,
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
	err := service.CatChangeStatus(ctx, id, &status)

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

	status := CatStatusStruct{
		Status: CatStatusSick,
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, id).
		Return(nil).Once()
	mockRepo.EXPECT().
		CatChangeStatus(ctx, id, &status.Status).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.CatChangeStatus(ctx, id, &status)

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
	mockRepo.EXPECT().CatDelete(ctx, id).
		Return(nil).Once()

	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.CatDelete(ctx, id)

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
	err := service.CatDelete(ctx, id)

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
		CatDelete(ctx, id).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.CatDelete(ctx, id)

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
	input := &ToyCreate{
		Name:     "Gilbert",
		Type:     "bear",
		Color:    "white",
		Material: "synthetic",
	}

	// Mock repository
	var captured *ToyCreateRepo
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToyCreate(ctx, catID, mock.Anything).
		Run(func(ctx context.Context, catID *CatID, data *ToyCreateRepo) {
			captured = data
		}).
		Return(nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toyID, err := service.ToyCreate(ctx, catID, input)

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
		assert.Positive(t, captured.CreatedAt)
		assert.Positive(t, captured.UpdatedAt)
	}

	mockRepo.AssertExpectations(t)
}

func TestCatService_AddToy_CatExists_RepoErr(t *testing.T) {
	ctx := context.Background()

	// Input data
	id := &CatID{
		ID: "00000000-0000-0000-0000-000000000000",
	}

	input := &ToyCreate{
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
	_, err := service.ToyCreate(ctx, id, input)

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
	input := &ToyCreate{
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
		ToyCreate(ctx, catID, mock.Anything).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toyID, err := service.ToyCreate(ctx, catID, input)

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

	// input data
	name := "Tom"
	color := "green"
	params := &ToySearchParams{
		Name:  &name,
		Color: &color,
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
		ToySearch(ctx, catID, params).
		Return(exceptedToys, nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toys, err := service.ToySearch(ctx, catID, params)

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

	// input data
	name := "Tom"
	color := "green"
	params := &ToySearchParams{
		Name:  &name,
		Color: &color,
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
	_, err := service.ToySearch(ctx, id, params)

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

	// input data
	name := "Tom"
	color := "green"
	params := &ToySearchParams{
		Name:  &name,
		Color: &color,
	}

	// Mock repository
	mockRepo := NewMockCatRepositoryInterface(t)
	mockRepo.EXPECT().
		CatExists(ctx, catID).
		Return(nil).Once()
	mockRepo.EXPECT().
		ToySearch(ctx, catID, params).
		Return(nil, errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toyID, err := service.ToySearch(ctx, catID, params)

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
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
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
		ToyByID(ctx, toyID).
		Return(exceptedToy, nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toy, err := service.ToyByID(ctx, catID, toyID)

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
	toy, err := service.ToyByID(ctx, catID, toyID)

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
	toy, err := service.ToyByID(ctx, catID, toyID)

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
	toy, err := service.ToyByID(ctx, catID, toyID)

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
		ToyByID(ctx, toyID).
		Return(nil, errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	toy, err := service.ToyByID(ctx, catID, toyID)

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
		ToyDelete(ctx, toyID).
		Return(nil).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.ToyDelete(ctx, catID, toyID)

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
	err := service.ToyDelete(ctx, catID, toyID)

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
	err := service.ToyDelete(ctx, catID, toyID)

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
	err := service.ToyDelete(ctx, catID, toyID)

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
		ToyDelete(ctx, toyID).
		Return(errors.New("db error")).Once()

	// Service
	logger := zerolog.Nop()
	service := NewCatService(mockRepo, &logger)

	// Call method
	err := service.ToyDelete(ctx, catID, toyID)

	// Assertions
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
