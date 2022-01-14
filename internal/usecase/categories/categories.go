package categories

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/response"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/validation"
	"github.com/IndominusByte/magicimage"
	"github.com/gorilla/schema"
)

type CategoriesUsecase struct {
	categoriesRepo categoriesRepo
}

func NewCategoriesUsecase(categoryRepo categoriesRepo) *CategoriesUsecase {
	return &CategoriesUsecase{
		categoriesRepo: categoryRepo,
	}
}

type CreateUpdateSchema struct {
	// Id   int   `validate:"required"`
	name []string `validate:"required,min=1,dive,required"`
}

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var decoder = schema.NewDecoder()

func (uc *CategoriesUsecase) CreateCategory(ctx context.Context, rw http.ResponseWriter, file, payload interface{}) {
	magic := magicimage.New(file.(*multipart.Form))
	if err := magic.ValidateSingleImage("icon"); err != nil {
		response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
			"icon": err.Error(),
		})
		return
	}

	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	response.WriteJSONResponse(rw, 201, nil, map[string]interface{}{
		"_app": []string{
			"Successfully create category.",
		},
	})
}
