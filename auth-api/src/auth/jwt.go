package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/config"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/models"
	"time"
)

type CustomClaims struct {
	User models.User `json:"user"`
	jwt.RegisteredClaims
}

func newToken(user models.User) (string, error) {
	claims := CustomClaims{
		user,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.AllConfigs.SignKey))

}

func getAcl(token string) (user models.User) {

	tk, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AllConfigs.SignKey), nil
	})

	if err != nil {
		return models.User{}
	}

	if claim, ok := tk.Claims.(*CustomClaims); ok && tk.Valid {
		return claim.User
	} else {
		return models.User{}
	}

}
