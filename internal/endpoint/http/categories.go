package endpoint_http

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/IndominusByte/learn-go-restful-api/internal/constant"
	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/parser"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/response"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/validation"
	"github.com/go-chi/chi/v5"
)

type categoriesUsecaseIface interface {
	CreateCategory(ctx context.Context, rw http.ResponseWriter, file *multipart.Form, payload *categoriesentity.FormCreateUpdateSchema)
	GetAllCategory(ctx context.Context, rw http.ResponseWriter, payload *categoriesentity.QueryParamAllCategorySchema)
	GetCategoryById(ctx context.Context, rw http.ResponseWriter, categoryId int)
	UpdateCategory(ctx context.Context, rw http.ResponseWriter, file *multipart.Form, payload *categoriesentity.FormCreateUpdateSchema, categoryId int)
	DeleteCategoryById(ctx context.Context, rw http.ResponseWriter, categoryId int)
}

func AddCategories(r *chi.Mux, uc categoriesUsecaseIface) {
	r.Route("/categories", func(r chi.Router) {
		r.Post("/", func(rw http.ResponseWriter, r *http.Request) {
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": constant.FailedParseBody,
				})
				return
			}

			var p categoriesentity.FormCreateUpdateSchema

			if err := validation.ParseRequest(&p, r.Form); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": constant.FailedParseBody,
				})
				return
			}

			uc.CreateCategory(r.Context(), rw, r.MultipartForm, &p)
		})
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			var p categoriesentity.QueryParamAllCategorySchema

			if err := validation.ParseRequest(&p, r.URL.Query()); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": constant.FailedParseBody,
				})
				return
			}

			uc.GetAllCategory(r.Context(), rw, &p)
		})
		r.Get("/{category_id:[1-9][0-9]*}", func(rw http.ResponseWriter, r *http.Request) {
			categoryId, _ := parser.ParsePathToInt("/categories/(.*)", r.URL.Path)

			uc.GetCategoryById(r.Context(), rw, categoryId)
		})
		r.Put("/{category_id:[1-9][0-9]*}", func(rw http.ResponseWriter, r *http.Request) {
			categoryId, _ := parser.ParsePathToInt("/categories/(.*)", r.URL.Path)

			if err := r.ParseMultipartForm(32 << 20); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": constant.FailedParseBody,
				})
				return
			}

			var p categoriesentity.FormCreateUpdateSchema

			if err := validation.ParseRequest(&p, r.Form); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": constant.FailedParseBody,
				})
				return
			}

			uc.UpdateCategory(r.Context(), rw, r.MultipartForm, &p, categoryId)
		})
		r.Delete("/{category_id:[1-9][0-9]*}", func(rw http.ResponseWriter, r *http.Request) {
			categoryId, _ := parser.ParsePathToInt("/categories/(.*)", r.URL.Path)

			uc.DeleteCategoryById(r.Context(), rw, categoryId)
		})
	})
}
