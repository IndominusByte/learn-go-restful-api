package categories

import (
	"context"

	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
)

type categoriesRepo interface {
	GetCategoryByName(ctx context.Context, payload *categoriesentity.FormCreateSchema) (*categoriesentity.Category, error)
	InsertCategory(ctx context.Context, payload *categoriesentity.FormCreateSchema) int
	GetAllCategoryPaginate(ctx context.Context, payload *categoriesentity.QueryParamAllCategorySchema) (*categoriesentity.CategoryPaginate, error)
}
