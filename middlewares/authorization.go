package middlewares

import (
	"github.com/OhMinsSup/lafu-server/lib"
	"github.com/labstack/echo/v4"
	"log"
	"strings"
	"time"
)


func Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, accessErr := c.Cookie("access_token")
		refreshToken, refreshErr := c.Cookie("refresh_token")
		log.Println(accessErr, refreshErr)
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

		exp := accessTokenData["exp"].(float64)
		target := time.Unix(int64(exp), 0).AddDate(0, 0, -1)
		expireDate := time.Hour * 24
		diff := time.Since(target) > expireDate

		if diff && refreshToken != nil {
			log.Println("refreshToken 발급")
			// refresh(context, refreshToken.Value)
		}

		if refreshErr != nil {
			log.Println("refreshToken 없음")
			return next(c)
		}
		log.Println("미들웨어....")
		return next(c)
	}
}