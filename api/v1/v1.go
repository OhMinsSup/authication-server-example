package v1

import (
	"github.com/OhMinsSup/lafu-server/api/v1/auth"
	"github.com/OhMinsSup/lafu-server/api/v1/file"
	"github.com/OhMinsSup/lafu-server/api/v1/promotion"
	"github.com/labstack/echo/v4"
)

// ApplyRoutes 라우터
func ApplyRoutes(e *echo.Group) {
	v1 := e.Group("/v1.0")
	auth.ApplyRoutes(v1)
	file.ApplyRoutes(v1)
	promotion.ApplyRoutes(v1)
}
