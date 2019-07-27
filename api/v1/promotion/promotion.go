package promotion

import (
"github.com/labstack/echo/v4"
)

// ApplyRoutes 라우터
func ApplyRoutes(e *echo.Group) {
	promotion := e.Group("/promotion")

	promotion.POST("/sendEmail", promotionEmail)
}