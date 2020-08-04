package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Need to make an env file and store db connection
// and secret Key in.
var (
	SecretKey = []byte("secret")
)

// GenerateToken generates a jwt and assigns a username
// to it's claims and returns it.
func GenerateToken(username string) (string, error) {
	//create a new jwt
	token := jwt.New(jwt.SigningMethodHS256)
	//create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)
	//Set token claims
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating Key")
		return "", err
	}
	return tokenString, nil
}

// PareseToken patses a jwt and returns the username in it's claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
