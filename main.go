package main

import (
	"github.com/OhMinsSup/lafu-server/middlewares"
	"os"

	"github.com/OhMinsSup/lafu-server/api"
	"github.com/OhMinsSup/lafu-server/database"
	"github.com/OhMinsSup/lafu-server/lib"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	port := os.Getenv("PORT")
	db, _ := database.Initialize()

	defer db.Close()

	e.Validator = lib.NewValidator()
	e.Logger.SetLevel(log.INFO)
	e.Use(database.Inject(db))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middlewares.Authorization)

	api.ApplyRoutes(e)

	e.Logger.Fatal(e.Start(":" + port))
}
