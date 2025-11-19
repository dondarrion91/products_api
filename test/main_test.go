package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"project/cmd/routes"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUpRoute(t *testing.T) {
	// Crear Echo
	e := echo.New()

	// Registrar las rutas igual que en main.go
	router := routes.Routes(e)

	// Crear request GET /
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Ejecutar la request en el router
	router.ServeHTTP(rec, req)

	// Assertions
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "UP", rec.Body.String())
}
