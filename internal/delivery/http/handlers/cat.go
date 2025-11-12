package handlers

import (
	"context"

	"github.com/Irurnnen/go-gin-template/internal/services/cat"
	"github.com/gin-gonic/gin"
)

type (
	CatHandlerInterface interface {
		Create(ctx *gin.Context)
		Search(ctx *gin.Context)
		GetByID(ctx *gin.Context)
		ChangeStatus(ctx *gin.Context)
		Delete(ctx *gin.Context)
		AddToy(ctx *gin.Context)
		GetToys(ctx *gin.Context)
		GetToyByID(ctx *gin.Context)
		DeleteToy(ctx *gin.Context)
	}

	CatServiceInterface interface {
		Create(ctx context.Context, data *cat.CreateCat) (*cat.CatID, error)
		Search(ctx context.Context, params *cat.CatSearchParams) ([]*cat.Cat, error)
		GetByID(ctx context.Context, id *cat.CatID) (*cat.Cat, error)
		ChangeStatus(ctx context.Context, id *cat.CatID, status *cat.CatStatus) error
		Delete(ctx context.Context, id *cat.CatID) error
		AddToy(ctx context.Context, catID *cat.CatID, data *cat.CreateToy) (*cat.ToyID, error)
		GetToys(ctx context.Context, catID *cat.CatID) ([]*cat.Toy, error)
		GetToyByID(ctx context.Context, catID *cat.CatID, toyID *cat.ToyID) (*cat.Toy, error)
		DeleteToy(ctx context.Context, catID *cat.CatID, toyID *cat.ToyID) error
	}
)
