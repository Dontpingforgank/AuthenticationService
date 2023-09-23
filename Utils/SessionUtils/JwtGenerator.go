package SessionUtils

import (
	"errors"
	"github.com/Dontpingforgank/AuthenticationService/Models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJwtToken(key string, claimsModel Models.UserClaimsModel) (string, error) {
	if key != "" && claimsModel.Id > 0 {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  claimsModel.Id,
			"exp": time.Now().Add(time.Hour * 24 * 16).Unix(),
		})

		tokenString, err := token.SignedString([]byte(key))
		if err != nil {
			return "", err
		}

		return tokenString, nil
	}

	return "", errors.New("can't generate session token")
}
