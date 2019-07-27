package lib

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
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
func GenerateUserToken(accessData JSON, refreshData JSON) (map[string]string, error) {
	accessDate := time.Now().Add(time.Hour * 24 * 7)
	refreshDate := time.Now().Add(time.Hour * 24 * 30)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": accessData,
		"exp": accessDate.Unix(),
	})
	log.Println(accessToken)
	access, errAccess := accessToken.SignedString(jwtKey)
	log.Println("access: ",errAccess)
	if errAccess != nil {
		return nil, errAccess
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": refreshData,
		"exp": refreshDate.Unix(),
	})
	refresh, errRefresh := refreshToken.SignedString(jwtKey)
	log.Println("refresh: ",errRefresh)
	if errRefresh != nil {
		return nil, errRefresh
	}

	return map[string]string{
		"access_token": access,
		"refresh_token": refresh,
	}, nil
}

// RefreshUserToken 토큰을 생성하는 함수
func RefreshUserToken(accessData JSON,  refreshData JSON, refreshTokenExp int64, originalRefreshToken string) (map[string]string, error) {
	var refresh = originalRefreshToken
	target := time.Unix(refreshTokenExp, 0).AddDate(0, 0, 30)
	expireDate := time.Hour * 24 * 30
	diff := time.Since(target) > expireDate
	if diff {
		log.Println("...refreshing refreshToken")
		refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": refreshData,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		 })

		refreshResult, err := refreshToken.SignedString(jwtKey)
		if err != nil {
			return nil, err
		}
		refresh = refreshResult
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": accessData,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	access, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token": access,
		"refresh_token": refresh,
	}, nil
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