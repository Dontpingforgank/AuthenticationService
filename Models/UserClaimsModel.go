package Models

import "github.com/golang-jwt/jwt/v5"

type UserClaimsModel struct {
	Id  int   `json:"id"`
	Exp int64 `json:"exp"`
	jwt.RegisteredClaims
}
