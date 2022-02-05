package endpoint_http

import (
	"net/http"

	"github.com/IndominusByte/learn-go-restful-api/internal/config"
	"github.com/creent-production/cdk-go/auth"
	"github.com/creent-production/cdk-go/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
)

func AddToken(r *chi.Mux, redisCli *redis.Pool, cfg *config.Config) {
	r.Route("/token", func(r chi.Router) {
		r.Get("/get-token", func(rw http.ResponseWriter, r *http.Request) {
			accessToken := auth.GenerateAccessToken(&auth.AccessToken{Sub: "1", Exp: jwtauth.ExpireIn(cfg.JWT.AccessExpires), Fresh: true})
			refreshToken := auth.GenerateRefreshToken(&auth.RefreshToken{Sub: "1", Exp: jwtauth.ExpireIn(cfg.JWT.RefreshExpires)})

			response.WriteJSONResponse(rw, 200, map[string]interface{}{
				// "access_token":  auth.NewJwtTokenHS([]byte(cfg.JWT.SecretKey), cfg.JWT.Algorithm, accessToken),
				// "refresh_token": auth.NewJwtTokenHS([]byte(cfg.JWT.SecretKey), cfg.JWT.Algorithm, refreshToken),
				"access_token":  auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, accessToken),
				"refresh_token": auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, refreshToken),
			}, nil)
		})

		r.Group(func(r chi.Router) {
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

			// access-revoke
			r.Delete("/access-revoke", func(rw http.ResponseWriter, r *http.Request) {
				conn := redisCli.Get()
				defer conn.Close()

				_, claims, _ := jwtauth.FromContext(r.Context())
				conn.Do("SET", claims["jti"], "ok")

				response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
					"_app": "An access token has revoked.",
				})
			})
		})
		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					if err := auth.ValidateJWT(r.Context(), redisCli, "jwtRefreshRequired"); err != nil {
						response.WriteJSONResponse(rw, 401, nil, map[string]interface{}{
							"_header": err.Error(),
						})
						return
					}
					// Token is authenticated, pass it through
					next.ServeHTTP(rw, r)
				})
			})

			// refresh-token
			r.Post("/refresh-token", func(rw http.ResponseWriter, r *http.Request) {
				accessToken := auth.GenerateAccessToken(&auth.AccessToken{Sub: "1", Exp: jwtauth.ExpireIn(cfg.JWT.AccessExpires), Fresh: false})

				response.WriteJSONResponse(rw, 200, map[string]interface{}{
					// "access_token": auth.NewJwtTokenHS([]byte(cfg.JWT.SecretKey), cfg.JWT.Algorithm, accessToken),
					"access_token": auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, accessToken),
				}, nil)
			})
			// refresh-revoke
			r.Delete("/refresh-revoke", func(rw http.ResponseWriter, r *http.Request) {
				conn := redisCli.Get()
				defer conn.Close()

				_, claims, _ := jwtauth.FromContext(r.Context())
				conn.Do("SET", claims["jti"], "ok")

				response.WriteJSONResponse(rw, 200, nil, map[string]interface{}{
					"_app": "An refresh token has revoked.",
				})
			})
		})
	})
}
