package categories

import (
	"context"

	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
)

type categoriesRepo interface {
	GetCategoryByName(ctx context.Context, payload *categoriesentity.FormCreateUpdateSchema) (*categoriesentity.Category, error)
	InsertCategory(ctx context.Context, payload *categoriesentity.FormCreateUpdateSchema) int
	UpdateCategory(ctx context.Context, payload *categoriesentity.FormCreateUpdateSchema) error
	GetAllCategoryPaginate(ctx context.Context, payload *categoriesentity.QueryParamAllCategorySchema) (*categoriesentity.CategoryPaginate, error)
	GetCategoryById(ctx context.Context, categoryId int) (*categoriesentity.Category, error)
	DeleteCategoryById(ctx context.Context, categoryId int) error
}
