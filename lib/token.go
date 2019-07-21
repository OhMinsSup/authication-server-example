package lib

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"os"
	"time"
)

var jwtKey []byte

func init() {
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		panic("Failed to load secret key file")
	}
	jwtKey = key
}

// GenerateToken 토큰을 생성하는 함수
func Generate(data JSON) (string, error) {
	date := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  date.Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

// DecodeToken 토큰으로 생성된 데이터를 해석해서 반환하는 함수
func Decode(deocedToken string) (JSON, error) {
	result, err := jwt.Parse(deocedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return JSON{}, err
	}

	if !result.Valid {
		return JSON{}, errors.New("invalid token")
	}

	return result.Claims.(jwt.MapClaims), nil
}