package handlers

import (
	"context"
	"net/http"

	delHTTP "github.com/Irurnnen/go-gin-template/internal/delivery/http"
	"github.com/Irurnnen/go-gin-template/internal/delivery/http/dto"
	"github.com/Irurnnen/go-gin-template/internal/delivery/http/dto/request"
	"github.com/Irurnnen/go-gin-template/internal/delivery/http/dto/response"
	"github.com/Irurnnen/go-gin-template/internal/services/cat"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type (
	CatServiceInterface interface {
		CatCreate(ctx context.Context, data *cat.CatCreate) (*cat.CatID, error)
		CatSearch(ctx context.Context, params *cat.CatSearchParams) ([]*cat.Cat, error)
		CatByID(ctx context.Context, id *cat.CatID) (*cat.Cat, error)
		CatChangeStatus(ctx context.Context, id *cat.CatID, status *cat.CatStatusStruct) error
		CatDelete(ctx context.Context, id *cat.CatID) error
		ToyCreate(ctx context.Context, catID *cat.CatID, data *cat.ToyCreate) (*cat.ToyID, error)
		ToySearch(ctx context.Context, catID *cat.CatID, params *cat.ToySearchParams) ([]*cat.Toy, error)
		ToyByID(ctx context.Context, catID *cat.CatID, toyID *cat.ToyID) (*cat.Toy, error)
		ToyDelete(ctx context.Context, catID *cat.CatID, toyID *cat.ToyID) error
	}

	CatHandler struct {
		service CatServiceInterface
		logger  *zerolog.Logger
	}
)

func NewCatHandler(service CatServiceInterface, logger *zerolog.Logger) *CatHandler {
	return &CatHandler{
		service: service,
		logger:  logger,
	}
}

// CatCreate godoc
//
//	@Summary		CatCreate new cat
//	@Description	create cat
//	@Tags			Cat
//	@Produce		json
//	@Param			CatCreate	body		request.CatCreate	true	"Data for creating a cat"
//	@Success		201			{object}	response.CatID
//	@Failure		400			{object}	dto.HTTPError
//	@Failure		500			{object}	dto.HTTPError
//	@Router			/cat [POST]
func (ch *CatHandler) CatCreate(ctx *gin.Context) {
	// Get data from the body
	createForm := new(request.CatCreate)
	err := ctx.ShouldBindBodyWithJSON(createForm)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingJSON, Code: http.StatusBadRequest})
		return
	}

	// Create cat
	id, err := ch.service.CatCreate(ctx.Request.Context(), &cat.CatCreate{
		Name:           createForm.Name,
		Breed:          createForm.Breed,
		BirthTimestamp: createForm.BirthTimestamp,
		Color:          createForm.Color,
		Gender:         createForm.Gender,
		Weight:         createForm.Weight,
		Description:    createForm.Description,
	})

	// Process errors
	if err != nil {
		// TODO: wrap into local error
		ctx.Error(err)
	}

	// Transform to dto
	catID := response.CatID{
		ID: id.ID,
	}

	// Send response
	ctx.JSON(http.StatusCreated, catID)
}

// CatSearch godoc
//
//	@Summary		CatSearch cats by filters
//	@Description	search cats
//	@Tags			Cat
//	@Produce		json
//	@Param			CatSearch	query		request.CatSearch	true	"Data for searching a cat"
//	@Success		200			{array}	response.Cat
//	@Success		400			{object}	dto.HTTPError
//	@Failure		500			{object}	dto.HTTPError
//	@Router			/cat [GET]
func (ch *CatHandler) CatSearch(ctx *gin.Context) {
	// Get data from the query
	searchParams := new(request.CatSearchParams)
	err := ctx.ShouldBindQuery(searchParams)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingQuery, Code: http.StatusBadRequest})
		return
	}

	// Search cats
	cats, err := ch.service.CatSearch(ctx.Request.Context(), &cat.CatSearchParams{
		Name:   searchParams.Name,
		Breed:  searchParams.Breed,
		Age:    searchParams.Age,
		Limit:  searchParams.Limit,
		Offset: searchParams.Offset,
	})

	// Process errors
	if err != nil {
		// TODO: wrap into local error
		ctx.Error(err)
	}

	// Transform to dto
	dtoCats := make([]*response.Cat, len(cats))
	for i, cat := range cats {
		dtoCats[i] = &response.Cat{
			ID:          cat.ID,
			Name:        cat.Name,
			Status:      string(cat.Status),
			Color:       cat.Color,
			Gender:      cat.Gender,
			Weight:      cat.Weight,
			Description: cat.Description,
			Breed:       cat.Breed,
			Age:         cat.Age,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		}
	}

	// Send response
	ctx.JSON(http.StatusCreated, dtoCats)
}

// CatByID godoc
//
//	@Summary		Get cat by its id
//	@Description	get by id
//	@Tags			Cat
//	@Produce		json
//	@Param			CatID	uri		request.CatID	true	"cats ID"
//	@Success		200			{object}	response.Cat
//	@Failure		400			{object}	dto.HTTPError
//	@Failure		404			{object}	dto.HTTPError
//	@Failure		500			{object}	dto.HTTPError
//	@Router			/cat/{cid} [GET]
func (ch *CatHandler) CatByID(ctx *gin.Context) {
	// Get data from the body
	catID := new(request.CatID)
	err := ctx.ShouldBindUri(catID)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingURI, Code: http.StatusBadRequest})
		return
	}

	// Get cat by id
	cat, err := ch.service.CatByID(ctx.Request.Context(), &cat.CatID{
		ID: catID.ID,
	})

	// Process errors
	if err != nil {
		// TODO: wrap into local error
		ctx.Error(err)
	}

	// Transform to dto
	dtoCat := response.Cat{
		ID:          cat.ID,
		Name:        cat.Name,
		Status:      string(cat.Status),
		Color:       cat.Color,
		Gender:      cat.Gender,
		Weight:      cat.Weight,
		Description: cat.Description,
		Breed:       cat.Breed,
		Age:         cat.Age,
		CreatedAt:   cat.CreatedAt,
		UpdatedAt:   cat.UpdatedAt,
	}

	// Send response
	ctx.JSON(http.StatusCreated, dtoCat)
}

// CatChangeStatus godoc
//
//	@Summary		Change cat status by its id
//	@Description	change cat status
//	@Tags			Cat
//	@Produce		json
//	@Param			CatID	uri		request.CatID	true	"cats ID"
//	@Param			CatStatus	body		request.CatStatus	true	"cats status"
//	@Success		200			{object}	dto.Message
//	@Failure		400			{object}	dto.HTTPError
//	@Failure		404			{object}	dto.HTTPError
//	@Failure		500			{object}	dto.HTTPError
//	@Router			/cat/{cid}/status [PATCH]
func (ch *CatHandler) CatChangeStatus(ctx *gin.Context) {
	// Get data from the body
	catID := new(request.CatID)
	err := ctx.ShouldBindUri(catID)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingURI, Code: http.StatusBadRequest})
		return
	}

	// Get data from the body
	status := new(request.CatStatus)
	err = ctx.ShouldBindJSON(status)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingJSON, Code: http.StatusBadRequest})
		return
	}

	// Change cat status
	err = ch.service.CatChangeStatus(
		ctx.Request.Context(),
		&cat.CatID{
			ID: catID.ID,
		},
		&cat.CatStatusStruct{
			Status: cat.CatStatus(status.Status),
		})

	// Process errors
	if err != nil {
		// TODO: wrap into local error
		ctx.Error(err)
	}

	// Send response
	ctx.JSON(http.StatusCreated, dto.Message{Message: "Status changed successfully"})
}

// CatDelete godoc
//
//	@Summary		CatDelete cat by its id
//	@Description	delete cat
//	@Tags			Cat
//	@Produce		json
//	@Param			CatID	uri		request.CatID	true	"cats ID"
//	@Success		200			{object}	dto.Message
//	@Failure		400			{object}	dto.HTTPError
//	@Failure		404			{object}	dto.HTTPError
//	@Failure		500			{object}	dto.HTTPError
//	@Router			/cat/{cid} [DELETE]
func (ch *CatHandler) CatDelete(ctx *gin.Context) {
	// Get data from the body
	catID := new(request.CatID)
	err := ctx.ShouldBindUri(catID)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingURI, Code: http.StatusBadRequest})
		return
	}

	// Delete cat
	err = ch.service.CatDelete(
		ctx.Request.Context(),
		&cat.CatID{
			ID: catID.ID,
		})

	// Process errors
	if err != nil {
		// TODO: wrap into local error
		ctx.Error(err)
	}

	// Send response
	ctx.JSON(http.StatusCreated, dto.Message{Message: "Cat deleted successfully"})
}

// ToyCreate godoc
//
//	@Summary		Add toy for cat by cat id
//	@Description	add toy
//	@Tags			Toy
//	@Produce		json
//	@Param			CatID	uri		request.CatID	true	"cats ID"
//	@Param			ToyCreate	json		request.ToyCreate	true	"Data for creating a toy"
//	@Success		201			{object}	response.ToyID
//	@Failure		400			{object}	dto.HTTPError
//	@Failure		404			{object}	dto.HTTPError
//	@Failure		500			{object}	dto.HTTPError
//	@Router			/cat/{cid}/toy [POST]
func (ch *CatHandler) ToyCreate(ctx *gin.Context) {
	// Get data from the body
	catID := new(request.CatID)
	err := ctx.ShouldBindUri(catID)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingURI, Code: http.StatusBadRequest})
		return
	}

	// Get data from the body
	createParams := new(request.ToyCreate)
	err = ctx.ShouldBindJSON(createParams)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingJSON, Code: http.StatusBadRequest})
		return
	}

	// Delete cat
	toyID, err := ch.service.ToyCreate(
		ctx.Request.Context(),
		&cat.CatID{
			ID: catID.ID,
		},
		&cat.ToyCreate{
			Name:     createParams.Name,
			Type:     createParams.Type,
			Color:    createParams.Color,
			Material: createParams.Material,
		})

	// Process errors
	if err != nil {
		// TODO: wrap into local error
		ctx.Error(err)
	}

	// Transform to dto
	dtoID := &response.ToyID{
		ID: toyID.ID,
	}

	// Send response
	ctx.JSON(http.StatusCreated, dtoID)
}

// ToySearch godoc
//
//	@Summary		Add toy for cat by cat id
//	@Description	add toy
//	@Tags			Toy
//	@Produce		json
//	@Param			CatID	uri		request.CatID	true	"cats ID"
//	@Param			ToyCreate	json		request.ToyCreate	true	"Data for creating a toy"
//	@Success		201			{object}	response.ToyID
//	@Failure		400			{object}	dto.HTTPError
//	@Failure		404			{object}	dto.HTTPError
//	@Failure		500			{object}	dto.HTTPError
//	@Router			/cat/{cid}/toy [POST]
func (ch *CatHandler) ToySearch(ctx *gin.Context) {
	// Get data from the body
	catID := new(request.CatID)
	err := ctx.ShouldBindUri(catID)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingURI, Code: http.StatusBadRequest})
		return
	}

	// Get data from the body
	searchParams := new(request.ToySearch)
	err = ctx.ShouldBindJSON(searchParams)
	if err != nil {
		ctx.Error(&delHTTP.BindingError{Err: err, Zone: delHTTP.BindingJSON, Code: http.StatusBadRequest})
		return
	}

	// Delete cat
	toys, err := ch.service.ToySearch(
		ctx.Request.Context(),
		&cat.CatID{
			ID: catID.ID,
		},
		&cat.ToySearchParams{
			Name:     searchParams.Name,
			Type:     searchParams.Type,
			Color:    searchParams.Color,
			Material: searchParams.Material,
			Limit:    searchParams.Limit,
			Offset:   searchParams.Offset,
		})

	// Process errors
	if err != nil {
		// TODO: wrap into local error
		ctx.Error(err)
	}

	// Transform to dto
	dtoToys := make([]*response.Toy, len(toys))
	for i, toy := range toys {
		dtoToys[i] = &response.Toy{
			ID:        toy.ID,
			Name:      toy.Name,
			Type:      toy.Type,
			Color:     toy.Color,
			Material:  toy.Material,
			CreatedAt: toy.CreatedAt,
			UpdatedAt: toy.UpdatedAt,
		}
	}
	// Send response
	ctx.JSON(http.StatusCreated, dtoToys)
}
