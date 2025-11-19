package main_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"project/internal/item_detail/utils"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------------------------------
// 🔧 Helper para crear un echo.Context fácilmente
// -------------------------------------------------------------
func newContext(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// -------------------------------------------------------------
// 🧪 TEST 1 — UnmarshalTypeError (invalid field type)
// -------------------------------------------------------------
func TestValidateBody_UnmarshalTypeError(t *testing.T) {
	t.Log("🔍 TEST: ValidateBody returns readable error for type mismatch")

	c, rec := newContext("POST", "/test", `{"images": "not-array"}`)

	// Simular error de tipo
	err := &json.UnmarshalTypeError{
		Field: "images",
		Type:  reflect.TypeOf([]string{}),
	}

	_ = utils.ValidateBody(c, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid field 'images'")

	t.Log("✅ Returned proper error for invalid 'images' field type")
}

// -------------------------------------------------------------
// 🧪 TEST 2 — UnmarshalTypeError genérico
// -------------------------------------------------------------
func TestValidateBody_UnmarshalTypeError_Generic(t *testing.T) {
	t.Log("🔍 TEST: ValidateBody returns generic invalid type error")

	c, rec := newContext("POST", "/test", `{"price": "not-number"}`)

	err := &json.UnmarshalTypeError{
		Field: "price",
		Type:  reflect.TypeOf(0.0),
	}

	_ = utils.ValidateBody(c, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"field":"price"`)

	t.Log("✅ Returned proper message for generic UnmarshalTypeError")
}

// -------------------------------------------------------------
// 🧪 TEST 3 — Validator errors
// -------------------------------------------------------------
type dummyPayload struct {
	Title string `validate:"required"`
}

func TestValidateBody_ValidatorErrors(t *testing.T) {
	t.Log("🔍 TEST: ValidateBody transforms validator errors into readable messages")

	c, rec := newContext("POST", "/test", `{"title":""}`)

	validate := validator.New()
	var payload dummyPayload
	_ = c.Bind(&payload)

	err := validate.Struct(payload)

	_ = utils.ValidateBody(c, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "validation error")
	assert.Contains(t, rec.Body.String(), "Title")

	t.Log("✅ Properly formatted validator error response")
}

// -------------------------------------------------------------
// 🧪 TEST 4 — Empty JSON / EOF
// -------------------------------------------------------------
func TestValidateBody_EmptyBody(t *testing.T) {
	t.Log("🔍 TEST: ValidateBody detects empty or malformed JSON body (EOF)")

	c, rec := newContext("POST", "/test", ``)

	// Simular EOF
	err := errors.New("EOF")

	_ = utils.ValidateBody(c, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "empty or malformed JSON body")

	t.Log("✅ Empty/malformed JSON body detected")
}

// -------------------------------------------------------------
// 🧪 TEST 5 — Generic error
// -------------------------------------------------------------
func TestValidateBody_GenericError(t *testing.T) {
	t.Log("🔍 TEST: ValidateBody returns generic invalid request body")

	c, rec := newContext("POST", "/test", `{}`)

	err := errors.New("random error")

	_ = utils.ValidateBody(c, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid request body")

	t.Log("✅ Returned generic invalid body error as expected")
}
