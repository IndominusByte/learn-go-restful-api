package categories

import (
	"context"
	"mime/multipart"
	"net/http"

	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/response"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/validation"
	"github.com/IndominusByte/magicimage"
)

type CategoriesUsecase struct {
	categoriesRepo categoriesRepo
}

func NewCategoriesUsecase(categoryRepo categoriesRepo) *CategoriesUsecase {
	return &CategoriesUsecase{
		categoriesRepo: categoryRepo,
	}
}

func (uc *CategoriesUsecase) CreateCategory(ctx context.Context, rw http.ResponseWriter,
	file *multipart.Form, payload *categoriesentity.FormCreateSchema) {

	magic := magicimage.New(file)
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

	if _, err := uc.categoriesRepo.GetCategoryByName(ctx, payload); err == nil {
		response.WriteJSONResponse(rw, 400, nil, map[string]interface{}{
			"_app": "The name has already been taken.",
		})
		return
	}

	magic.SaveImages(100, 100, "static/icon-categories", true)
	payload.Icon = magic.FileNames[0]

	// save into database
	uc.categoriesRepo.InsertCategory(ctx, payload)

	response.WriteJSONResponse(rw, 201, nil, map[string]interface{}{
		"_app": "Successfully create category.",
	})
}

func (uc *CategoriesUsecase) GetAllCategory(ctx context.Context, rw http.ResponseWriter,
	payload *categoriesentity.QueryParamAllCategorySchema) {

	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	t, _ := uc.categoriesRepo.GetAllCategoryPaginate(ctx, payload)

	response.WriteJSONResponse(rw, 200, t, nil)
}
