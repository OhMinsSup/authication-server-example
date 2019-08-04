package middlewares

import (
	"github.com/OhMinsSup/lafu-server/database/models"
	"github.com/OhMinsSup/lafu-server/lib"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	_ "net/http"
	"strings"
	"time"
)

func refresh(c echo.Context, refreshToken, accessToken string) (string, error) {
	db := c.Get("db").(*gorm.DB)
	refreshDecoded, errR := lib.Decode(refreshToken)
	if errR != nil {
		log.Println(errR)
		return "", errR
	}

	accessDecoded, errA := lib.Decode(accessToken)
	if errA != nil {
		log.Println(errA)
		return "", errA
	}

	refreshData := refreshDecoded["user"].(lib.JSON)

	var user models.User
	if errU := db.Where("id = ?", refreshData["id"]).First(&user).Error; errU != nil {
		log.Println(errU)
		return "", errU
	}

	exp := refreshDecoded["exp"].(int64)

	tokens, errT := lib.RefreshUserToken(accessDecoded, refreshDecoded, exp , refreshToken)
	if errT != nil {
		log.Println(errT)
		return "", errT
	}

	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = "access_token"
	accessTokenCookie.Value = tokens["access_token"]
	accessTokenCookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	accessTokenCookie.Path= "/"
	accessTokenCookie.HttpOnly = true
	accessTokenCookie.Secure = false

	refreshTokenCookie := new (http.Cookie)
	refreshTokenCookie.Name = "refresh_token"
	refreshTokenCookie.Value = tokens["refresh_token"]
	refreshTokenCookie.Expires = time.Now().Add(time.Hour * 24 * 30)
	refreshTokenCookie.Path = "/"
	refreshTokenCookie.HttpOnly = true
	refreshTokenCookie.Secure = false

	c.SetCookie(accessTokenCookie)
	c.SetCookie(refreshTokenCookie)
	return refreshData["id"].(string), nil
}

func Authorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		exists := c.Get("userId")
		if exists != nil {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "로그인후 이용해주세요")
	}
}

func Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, accessErr := c.Cookie("access_token")
		refreshToken, refreshErr := c.Cookie("refresh_token")
		if accessErr != nil {
			log.Println("Authorization...")
			authorization := c.Request().Header.Get("Authorization")
			if authorization == "" {
				log.Println("Authorization 없음(1)")
				return next(c)
			}

			sp := strings.Split(authorization, "Bearer ")
			if len(sp) < 1 {
				log.Println("Authorization 없음(2)")
				return next(c)
			}
			accessToken.Value = sp[1]
		}

		accessTokenData, dataErr := lib.Decode(accessToken.Value)
		if dataErr != nil {
			log.Println("Decoded 토큰")
			return next(c)
		}

		accessData := accessTokenData["user"].(lib.JSON)
		c.Set("userId", accessData["id"])

		exp := accessTokenData["exp"].(float64)
		target := time.Unix(int64(exp), 0).AddDate(0, 0, -1)
		expireDate := time.Hour * 24
		diff := time.Since(target) > expireDate
		if diff && refreshToken != nil {
			log.Println("refreshToken 발급")
			userId, err := refresh(c, refreshToken.Value, accessToken.Value)
			if err != nil {
				return next(c)
			}
			c.Set("userId", userId)
		}

		if refreshErr != nil {
			log.Println("refreshToken 없음")
			return next(c)
		}
		log.Println("authorization")
		return next(c)
	}
}