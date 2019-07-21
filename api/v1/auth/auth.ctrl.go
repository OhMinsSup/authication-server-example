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

	serialized := user.Serialize()
	token, _ := lib.Generate(serialized)

	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
		"user": serialized,
		"token": token,
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