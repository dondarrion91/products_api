package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ---------------------
// Custom errors
// ---------------------
var (
	ErrInvalidImages = errors.New("invalid field 'images'")
	ErrEmptyBody     = errors.New("empty or malformed JSON body")
	ErrInvalidBody   = errors.New("invalid request body")
)

// ---------------------
// ValidateBody
// ---------------------
func ValidateBody(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	// Caso 1️⃣: Error de tipo al hacer Bind (por ejemplo, string donde se esperaba array)
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		if unmarshalErr.Field == "images" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": ErrInvalidImages.Error(),
				"hint":  "must be an array, e.g. [{\"url\": \"...\"}]",
			})
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "invalid field type",
			"field":  unmarshalErr.Field,
			"detail": err.Error(),
		})
	}

	// Caso 2️⃣: Error del validador (validator.v10)
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var fields []string
		for _, e := range validationErrs {
			// Generar mensajes legibles según el tag de validación
			var msg string
			var prop = strings.Fields(e.Error())[1]

			switch e.Tag() {
			case "required":
				msg = prop + " is required"
			case "min":
				msg = fmt.Sprintf("%v must have at least %s characters", prop, e.Param())
			case "max":
				msg = fmt.Sprintf("%v must have at most %s characters", prop, e.Param())
			case "url":
				msg = "must be a valid URL"
			default:
				msg = fmt.Sprintf("failed rule '%s'", e.Tag())
			}

			fields = append(fields, fmt.Sprintf("%s (%s)", e.Field(), msg))
		}

		return c.JSON(http.StatusBadRequest, echo.Map{
			"error":  "validation error",
			"fields": fields,
		})
	}

	// Caso 3️⃣: JSON vacío o mal formado
	if strings.Contains(err.Error(), "EOF") {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": ErrEmptyBody.Error(),
		})
	}

	// Caso 4️⃣: Error genérico
	return c.JSON(http.StatusBadRequest, echo.Map{
		"error": ErrInvalidBody.Error(),
	})
}
