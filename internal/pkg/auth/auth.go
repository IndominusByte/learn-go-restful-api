package auth

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-chi/jwtauth"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

type RefreshToken struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

type AccessToken struct {
	Sub   string `json:"sub"`
	Exp   int64  `json:"exp"`
	Fresh bool   `json:"fresh"`
}

func DecodeRSA(public, private string) (interface{}, interface{}) {
	privateKeyBlock, _ := pem.Decode([]byte(private))
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	publicKeyBlock, _ := pem.Decode([]byte(public))
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	return publicKey, privateKey
}

func NewJwtTokenHS(secret []byte, alg string, claims ...map[string]interface{}) string {
	token := jwt.New()
	if len(claims) > 0 {
		for k, v := range claims[0] {
			token.Set(k, v)
		}
	}
	tokenPayload, err := jwt.Sign(token, jwa.SignatureAlgorithm(alg), secret)
	if err != nil {
		panic(err)
	}
	return string(tokenPayload)
}

func NewJwtTokenRSA(public, private, alg string, claims map[string]interface{}) string {
	publicKey, privateKey := DecodeRSA(public, private)

	tokenAuth := jwtauth.New(alg, privateKey, publicKey)
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func GenerateAccessToken(at *AccessToken) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(at)
	json.Unmarshal(inrec, &inInterface)

	jwtauth.SetIssuedNow(inInterface)
	inInterface["jti"] = uuid.New().String()
	inInterface["type"] = "access"

	return inInterface
}

func GenerateRefreshToken(at *RefreshToken) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(at)
	json.Unmarshal(inrec, &inInterface)

	jwtauth.SetIssuedNow(inInterface)
	inInterface["jti"] = uuid.New().String()
	inInterface["type"] = "refresh"

	return inInterface
}

func ValidateJWT(ctx context.Context, redisCli *redis.Pool, typeJWT string) error {
	_, _, err := jwtauth.FromContext(ctx)

	switch typeJWT {
	case "jwtRequired", "jwtRefreshRequired", "jwtFreshRequired":
		return VerifyJWT(ctx, redisCli, typeJWT)
	case "jwtOptional":
		if err != nil && err.Error() == "no token found" {
			return nil
		}
		return VerifyJWT(ctx, redisCli, typeJWT)
	default:
		panic(fmt.Sprintf("Upps %s not found", typeJWT))
	}
}

func VerifyJWT(ctx context.Context, redisCli *redis.Pool, typeJWT string) error {
	token, claims, err := jwtauth.FromContext(ctx)
	// validate token
	if err != nil {
		return jwtauth.ErrorReason(err)
	}

	if err := jwt.Validate(token); err != nil {
		return jwtauth.ErrorReason(err)
	}

	// type token must be between 'access' or 'refresh'
	typeToken, ok := claims["type"]
	if !ok || (typeToken != "access" && typeToken != "refresh") {
		return errors.New("type token must be between access or refresh")
	}
	// if access token fresh must be included
	if _, ok := claims["fresh"]; !ok && typeToken == "access" {
		return errors.New("fresh must be included in token")
	}
	// check fresh is boolean
	if _, ok := claims["fresh"]; ok && reflect.TypeOf(claims["fresh"]).String() != "bool" {
		return errors.New("token is unauthorized")
	}
	// iat must be include
	if _, ok := claims["iat"]; !ok {
		return errors.New("iat must be included in token")
	}
	// jti must be include
	if _, ok := claims["jti"]; !ok {
		return errors.New("jti must be included in token")
	}
	// exp must be include
	if _, ok := claims["exp"]; !ok {
		return errors.New("exp must be included in token")
	}

	// validate type token
	switch typeJWT {
	case "jwtRequired", "jwtOptional":
		if claims["type"] != "access" {
			return errors.New("only access token are allowed")
		}
	case "jwtFreshRequired":
		if claims["type"] != "access" {
			return errors.New("only access token are allowed")
		}
		if claims["fresh"] == false {
			return errors.New("fresh token required")
		}
	case "jwtRefreshRequired":
		if claims["type"] != "refresh" {
			return errors.New("only refresh token are allowed")
		}
	}

	// check token is revoked
	conn := redisCli.Get()
	defer conn.Close()

	jti, _ := redis.String(conn.Do("GET", claims["jti"]))
	if len(jti) > 0 {
		return errors.New("token is revoked")
	}

	return nil
}
