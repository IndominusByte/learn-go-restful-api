package endpoint_http

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/IndominusByte/learn-go-restful-api/internal/constant"
	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/auth"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/parser"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/response"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/validation"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
)

type categoriesUsecaseIface interface {
	CreateCategory(ctx context.Context, rw http.ResponseWriter, file *multipart.Form, payload *categoriesentity.FormCreateUpdateSchema)
	GetAllCategory(ctx context.Context, rw http.ResponseWriter, payload *categoriesentity.QueryParamAllCategorySchema)
	GetCategoryById(ctx context.Context, rw http.ResponseWriter, categoryId int)
	UpdateCategory(ctx context.Context, rw http.ResponseWriter, file *multipart.Form, payload *categoriesentity.FormCreateUpdateSchema, categoryId int)
	DeleteCategoryById(ctx context.Context, rw http.ResponseWriter, categoryId int)
}

func AddCategories(r *chi.Mux, uc categoriesUsecaseIface, redisCli *redis.Pool) {
	r.Route("/categories", func(r chi.Router) {
		// Protected routes
		r.Group(func(r chi.Router) {
			// jwt token
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					if err := auth.ValidateJWT(r.Context(), redisCli, "jwtRequired"); err != nil {
						response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
							"_header": err.Error(),
						})
						return
					}
					// Token is authenticated, pass it through
					next.ServeHTTP(rw, r)
				})
			})

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

		// Public routes
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
	})
}
