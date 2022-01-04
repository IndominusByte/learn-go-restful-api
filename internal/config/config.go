package config

import (
	"fmt"
	"net/http"

	_ "github.com/IndominusByte/learn-go-restful-api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World bro!")
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), //The url pointing to API definition
	))
	http.ListenAndServe(":8080", r)
}
