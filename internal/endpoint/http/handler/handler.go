package handler_http

import (
	"net/http"
	"strings"

	"github.com/IndominusByte/learn-go-restful-api/internal/config"
	endpoint_http "github.com/IndominusByte/learn-go-restful-api/internal/endpoint/http"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/auth"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/mail"
	categoriesrepo "github.com/IndominusByte/learn-go-restful-api/internal/repo/categories"
	categoriesusecase "github.com/IndominusByte/learn-go-restful-api/internal/usecase/categories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
)

// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return fs.fs.Open("upps!")
	}

	return f, nil
}

type Server struct {
	Router *chi.Mux
	// Db config can be added here
	db       *sqlx.DB
	redisCli *redis.Pool
	cfg      *config.Config
}

func CreateNewServer(db *sqlx.DB, redisCli *redis.Pool, cfg *config.Config) *Server {
	s := &Server{db: db, redisCli: redisCli, cfg: cfg}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers() error {
	// jwt
	// TokenAuthHS256 := jwtauth.New(s.cfg.JWT.Algorithm, []byte(s.cfg.JWT.SecretKey), nil)
	// r.Use(jwtauth.Verifier(TokenAuthHS256))
	publicKey, privateKey := auth.DecodeRSA(s.cfg.JWT.PublicKey, s.cfg.JWT.PrivateKey)
	TokenAuthRS256 := jwtauth.New(s.cfg.JWT.Algorithm, privateKey, publicKey)
	s.Router.Use(jwtauth.Verifier(TokenAuthRS256))

	// middleware stack
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), //The url pointing to API definition
	))
	// serve file static
	fileServer := http.FileServer(FileSystem{http.Dir("static")})
	s.Router.Handle("/static/*", http.StripPrefix(strings.TrimRight("/static/", "/"), fileServer))

	// setup email
	m := mail.Mail{
		Server:   s.cfg.Mail.Server,
		Port:     s.cfg.Mail.Port,
		Username: s.cfg.Mail.Username,
		Password: s.cfg.Mail.Password,
	}

	// you can insert your behaviors here
	categoriesRepo, err := categoriesrepo.New(s.db)
	if err != nil {
		return err
	}
	categoriesUsecase := categoriesusecase.NewCategoriesUsecase(categoriesRepo)
	endpoint_http.AddCategories(s.Router, categoriesUsecase, s.redisCli)
	// add token
	endpoint_http.AddToken(s.Router, s.redisCli, s.cfg)
	// send email
	endpoint_http.AddEmail(s.Router, s.cfg, &m)

	return nil
}
