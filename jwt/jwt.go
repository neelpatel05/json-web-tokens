package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type User struct {
	Email string `json:"email"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var key = []byte("captainjacksparrow")

func GenerateJWT(user User) (string, int64) {

	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err!=nil {
		log.Fatal(err)
	}

	return tokenString, expirationTime.Unix()
}

func AuthorizeJWT(tokenString string) bool {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if !token.Valid {
		return false
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false
		}
		return false
	}
	return true

}

