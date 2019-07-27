package promotion

import (
	"fmt"
	"github.com/OhMinsSup/lafu-server/api/v1/auth/schema"
	"github.com/OhMinsSup/lafu-server/lib"
	"github.com/labstack/echo/v4"
	"net/http"
)

func promotionEmail(c echo.Context) error {
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