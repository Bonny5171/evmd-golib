package apicommon

import (
	"os"

	"github.com/go-chi/jwtauth"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	secret := os.Getenv("API_JWT_SECRET")
	if len(secret) == 0 {
		secret = "R3JlbmRlbmVAMjAxOQ"
	}

	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}
