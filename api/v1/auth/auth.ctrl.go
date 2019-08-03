package auth

import (
	"github.com/OhMinsSup/lafu-server/api/v1/auth/schema"
	"github.com/OhMinsSup/lafu-server/database/models"
	"github.com/OhMinsSup/lafu-server/lib"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func authInfo(c echo.Context) error {
	accessTokenCookie, err1 := c.Cookie("access_token")
	refreshTokenCookie, err2 := c.Cookie("refresh_token")
	if err1 != nil || err2 != nil {
		return c.JSON(404, echo.Map{
			"ok": false,
			"message": "토큰이 없습니다",
		})
	}

	return c.JSON(200, echo.Map{
		"ok": true,
		"userId": c.Get("userId"),
		"access_token": accessTokenCookie,
		"refresh_token": refreshTokenCookie,
	})
}

func localLogin(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	body := new (schema.LocalLoginSchema)

	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "WRONG_SCHEMA_BODY_DATA:"+err.Error())
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "VALIDATE_ERROR:"+err.Error())
	}

	var user models.User
	if err := db.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusForbidden, echo.Map{
			"ok": false,
			"message": "계정을 찾을 수 없습니다.",
		})
	}

	if !lib.Compare(body.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"ok": false,
			"message": "패스워드가 일치하지 않습니다",
		})
	}

	authToken := models.AuthToken{
		UserID: user.ID,
	}

	db.NewRecord(authToken)
	db.Create(&authToken)

	accessData := user.TokenData(nil)
	refreshData := user.TokenData(authToken.ID)
	serialized := user.Serialize()

	tokens, err := lib.GenerateUserToken(accessData, refreshData)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
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

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
		"user": serialized,
		"refreshToken": tokens["refresh_token"],
		"accessToken": tokens["access_token"],
	})
}

func localRegister(c echo.Context) error {
	db := c.Get("db").(*gorm.DB)
	body := new (schema.LocalRegisterSchema)

	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "WRONG_SCHEMA_BODY_DATA:"+err.Error())
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "VALIDATE_ERROR:"+err.Error())
	}

	var exists models.User
	if err := db.Where("username = ?", body.Username).Or("email = ?", body.Email).First(&exists).Error; err == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"ok": false,
			"message": "유저명또는 이메일이 이미 존재합니다.",
		})
	}

	hash, hashErr := lib.Hash(body.Password)
	if hashErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	user := models.User {
		Username: body.Username,
		Email: body.Email,
		Password: hash,
	}

	db.NewRecord(user)
	db.Create(&user)

	authToken := models.AuthToken{
		UserID:user.ID,
	}

	db.NewRecord(authToken)
	db.Create(&authToken)

	accessData := user.TokenData(nil)
	refreshData := user.TokenData(authToken.ID)
	serialized := user.Serialize()

	tokens, err := lib.GenerateUserToken(accessData, refreshData)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
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

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
		"user": serialized,
		"refreshToken": tokens["refresh_token"],
		"accessToken": tokens["access_token"],
	})
}

func logout(c echo.Context) error {
	access, errA := c.Cookie("access_token")
	if errA != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"ok": false,
			"message": errA,
		})
	}

	refresh, errR := c.Cookie("refresh_token")
	if errR != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"ok": false,
			"message": errR,
		})
	}

	refresh.Value = ""
	refresh.Path = "/"
	refresh.Expires = time.Now()
	refresh.HttpOnly = true
	refresh.Secure = false
	access.Value = ""
	access.Path = "/"
	access.Expires = time.Now()
	access.HttpOnly = true
	access.Secure = false

	c.SetCookie(access)
	c.SetCookie(refresh)
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}
