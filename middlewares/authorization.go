package middlewares

import "github.com/labstack/echo/v4"

func Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		accessToken, accessErr := context.Cookie("access_token")

		if accessErr != nil {

		}

		if err := next(context); err != nil {
			context.Error(err)
		}
	}
}