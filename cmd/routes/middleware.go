package routes

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func ApiKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	API_KEY := os.Getenv("API_KEY")

	return func(c echo.Context) error {

		key := c.Request().Header.Get("X-API-Key")

		// Si el header no está → bloquear
		if key == "" || API_KEY == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Missing X-API-Key header",
			})
		}

		// Si no coincide → bloquear
		if key != API_KEY {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "API key inválida",
			})
		}

		return next(c)
	}
}
