package tests

import (
	"net/http"
	"net/http/httptest"

	"github.com/IndominusByte/learn-go-restful-api/internal/config"
	handler_http "github.com/IndominusByte/learn-go-restful-api/internal/endpoint/http/handler"
)

func setupEnvironment() *handler_http.Server {
	// init config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	// connect the db
	db, err := config.DBConnect(cfg)
	if err != nil {
		panic(err)
	}
	// connect redis
	redisCli, err := config.RedisConnect(cfg)
	if err != nil {
		panic(err)
	}
	// mount router
	r := handler_http.CreateNewServer(db, redisCli, cfg)
	if err := r.MountHandlers(); err != nil {
		panic(err)
	}

	return r
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *handler_http.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}
