package endpoint_http

import (
	"net/http"
	"time"

	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/auth"
	"github.com/IndominusByte/learn-go-restful-api/internal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
)

var (
	PrivateKeyRS256String = `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBALxo3PCjFw4QjgOX06QCJIJBnXXNiEYwDLxxa5/7QyH6y77nCRQy
J3x3UwF9rUD0RCsp4sNdX5kOQ9PUyHyOtCUCAwEAAQJARjFLHtuj2zmPrwcBcjja
IS0Q3LKV8pA0LoCS+CdD+4QwCxeKFq0yEMZtMvcQOfqo9x9oAywFClMSlLRyl7ng
gQIhAOyerGbcdQxxwjwGpLS61Mprf4n2HzjwISg20cEEH1tfAiEAy9dXmgQpDPir
C6Q9QdLXpNgSB+o5CDqfor7TTyTCovsCIQDNCfpu795luDYN+dvD2JoIBfrwu9v2
ZO72f/pm/YGGlQIgUdRXyW9kH13wJFNBeBwxD27iBiVj0cbe8NFUONBUBmMCIQCN
jVK4eujt1lm/m60TlEhaWBC3p+3aPT2TqFPUigJ3RQ==
-----END RSA PRIVATE KEY-----
`
	PublicKeyRS256String = `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBALxo3PCjFw4QjgOX06QCJIJBnXXNiEYw
DLxxa5/7QyH6y77nCRQyJ3x3UwF9rUD0RCsp4sNdX5kOQ9PUyHyOtCUCAwEAAQ==
-----END PUBLIC KEY-----
`
)

func AddToken(r *chi.Mux, redisCli *redis.Pool) {
	r.Route("/token", func(r chi.Router) {
		r.Get("/get-token", func(rw http.ResponseWriter, r *http.Request) {
			accessToken := auth.GenerateAccessToken(&auth.AccessToken{Sub: "1", Exp: jwtauth.ExpireIn(5 * time.Minute), Fresh: true})
			refreshToken := auth.GenerateRefreshToken(&auth.RefreshToken{Sub: "1", Exp: jwtauth.ExpireIn(24 * time.Hour)})

			response.WriteJSONResponse(rw, 200, map[string]interface{}{
				"access_token":  auth.NewJwtTokenHS([]byte("secret"), "HS256", accessToken),
				"refresh_token": auth.NewJwtTokenHS([]byte("secret"), "HS256", refreshToken),
				// 	"access_token":  auth.NewJwtTokenRSA(PublicKeyRS256String, PrivateKeyRS256String, "RS256", accessToken),
				// 	"refresh_token": auth.NewJwtTokenRSA(PublicKeyRS256String, PrivateKeyRS256String, "RS256", refreshToken),
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
				accessToken := auth.GenerateAccessToken(&auth.AccessToken{Sub: "1", Exp: jwtauth.ExpireIn(5 * time.Minute), Fresh: true})

				response.WriteJSONResponse(rw, 200, map[string]interface{}{
					"access_token": auth.NewJwtTokenHS([]byte("secret"), "HS256", accessToken),
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
