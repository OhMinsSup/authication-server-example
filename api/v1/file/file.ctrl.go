package file

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func createUrl(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"ok": false,
			"message": "파일이 존재하지 않거나 올바른 형식의 파일이 아닙니다.",
		})
	}

	files := form.File["files"]
	log.Println(files)

	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}

func resize(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"ok": true,
	})
}