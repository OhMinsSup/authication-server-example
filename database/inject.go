package database

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

// Inject 데이터베이스 서버 미들웨어에 삽입
func Inject(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}
