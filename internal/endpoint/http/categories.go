package endpoint_http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IndominusByte/learn-go-restful-api/internal/entity/categories"
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
				fmt.Println("multipart")
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": "Invalid input type.",
				})
				return
			}

			var p categories.FormCreateSchema

			if err := validation.FormDecode(&p, r.Form); err != nil {
				fmt.Println("form")
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					"_body": "Invalid input type.",
				})
				return
			}

			if err := validation.StructValidate(p); err != nil {
				fmt.Println("struct")
				response.WriteJSONResponse(rw, 422, nil, err)
				return
			}

			fmt.Println(p, r.MultipartForm)
			uc.CreateCategory(r.Context(), rw, r.MultipartForm, &p)
		})
	})
}
