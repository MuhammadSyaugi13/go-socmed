package helper

import (
	"go-socmed/config"
	"go-socmed/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(user *entity.User) (string, error) {

	claims := JWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(config.ENV.JWT_KEY))

	return ss, err

}
