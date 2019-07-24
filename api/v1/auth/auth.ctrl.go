package auth

import (
	"fmt"
	"github.com/OhMinsSup/lafu-server/api/v1/auth/schema"
	"github.com/OhMinsSup/lafu-server/database/models"
	"github.com/OhMinsSup/lafu-server/lib"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

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
		return c.JSON(http.StatusConflict, echo.Map{
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

	accessData := user.TokenData("null")
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

	refreshTokenCookie := new (http.Cookie)
	refreshTokenCookie.Name = "refresh_token"
	refreshTokenCookie.Value = tokens["refresh_token"]
	refreshTokenCookie.Expires = time.Now().Add(time.Hour * 24 * 30)

	c.SetCookie(accessTokenCookie)
	c.SetCookie(refreshTokenCookie)

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
		"user": serialized,
		"refreshToken": tokens["refresh_token"],
		"accessToken": tokens["access_token"],
	})
}

func emailVerification(c echo.Context) error {
	// db := c.Get("db").(*gorm.DB)
	body := new (schema.EmailVerificationSchema)

	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "WRONG_SCHEMA_BODY_DATA:"+err.Error())
	}

	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "VALIDATE_ERROR:"+err.Error())
	}

	templateData := struct {
		keyword string
		Url string
	}{
		keyword: "키워드",
		Url: "https://www.naver.com",
	}

	m := lib.CreateSendEmail([]string{body.Email}, "키워드", "veloss<verification@gmail.com>")
	if err := m.ParseTemplate("statics/emailTemplate.html", templateData); err == nil {
		ok := m.SendEmail()
		fmt.Println(ok)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}