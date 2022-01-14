package endpoint_http

import (
	"context"
	"net/http"

	"github.com/IndominusByte/learn-go-restful-api/internal/constant"
	categoriesentity "github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/response"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/validation"
	"github.com/go-chi/chi/v5"
)

type categoriesUsecaseIface interface {
	CreateCategory(ctx context.Context, rw http.ResponseWriter, file, payload interface{})
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

			var p categoriesentity.FormCreateSchema

			if err := validation.FormDecode(&p, r.Form); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": constant.FailedParseBody,
				})
				return
			}

			uc.CreateCategory(r.Context(), rw, r.MultipartForm, &p)
		})
	})
}
