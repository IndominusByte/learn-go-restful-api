package categories

import (
	"context"
	"fmt"
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
	file *multipart.Form, payload *categoriesentity.FormCreateUpdateSchema) {

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

	magic.SaveImages(100, 100, "/app/static/icon-categories", true)
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

func (uc *CategoriesUsecase) GetCategoryById(ctx context.Context, rw http.ResponseWriter, categoryId int) {
	t, err := uc.categoriesRepo.GetCategoryById(ctx, categoryId)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			"_app": "Category not found.",
		})
		return
	}
	response.WriteJSONResponse(rw, 200, t, nil)
}

func (uc *CategoriesUsecase) UpdateCategory(ctx context.Context, rw http.ResponseWriter,
	file *multipart.Form, payload *categoriesentity.FormCreateUpdateSchema, categoryId int) {

	magic := magicimage.New(file)
	magic.Required = false
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

	category, err := uc.categoriesRepo.GetCategoryById(ctx, categoryId)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			"_app": "Category not found.",
		})
		return
	}

	if _, err := uc.categoriesRepo.GetCategoryByName(ctx, payload); err == nil && category.Name != payload.Name {
		response.WriteJSONResponse(rw, 400, nil, map[string]interface{}{
			"_app": "The name has already been taken.",
		})
		return
	}

	// delete the image from db if file exists
	if _, ok := file.File["icon"]; ok {
		magicimage.DeleteFolderAndFile(fmt.Sprintf("/app/static/icon-categories/%s", category.Icon))
		magic.SaveImages(100, 100, "/app/static/icon-categories", true)
		payload.Icon = magic.FileNames[0]
	}

	// update into db
	payload.Id = category.Id
	uc.categoriesRepo.UpdateCategory(ctx, payload)

	response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
		"_app": "Successfully update the category.",
	})
}

func (uc *CategoriesUsecase) DeleteCategoryById(ctx context.Context, rw http.ResponseWriter, categoryId int) {
	category, err := uc.categoriesRepo.GetCategoryById(ctx, categoryId)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			"_app": "Category not found.",
		})
		return
	}

	// delete into db
	magicimage.DeleteFolderAndFile(fmt.Sprintf("/app/static/icon-categories/%s", category.Icon))
	uc.categoriesRepo.DeleteCategoryById(ctx, category.Id)

	response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
		"_app": "Successfully delete the category.",
	})
}
