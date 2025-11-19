package main_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/item_detail/repo/datasource/dao"
	"project/internal/item_detail/rest"
	"project/internal/item_detail/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ---------- ENTIDAD MOCK ----------
type MockEntityHandler struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// ---------- MOCK DEL DAO (CrudDAO[T]) ----------
type MockCrudDAO struct {
	Data map[string]*MockEntityHandler
}

func NewMockCrudDAO() *MockCrudDAO {
	return &MockCrudDAO{
		Data: make(map[string]*MockEntityHandler),
	}
}

func (m *MockCrudDAO) Create(e *MockEntityHandler) (*MockEntityHandler, error) {
	m.Data[e.ID] = e
	return e, nil
}

func (m *MockCrudDAO) GetAll(q string, limit int, offset int) ([]*MockEntityHandler, error) {
	res := []*MockEntityHandler{}
	for _, v := range m.Data {
		res = append(res, v)
	}
	return res, nil
}

func (m *MockCrudDAO) GetByID(id string) (*MockEntityHandler, error) {
	v, ok := m.Data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (m *MockCrudDAO) Update(e *MockEntityHandler, id string) (*MockEntityHandler, error) {
	if _, ok := m.Data[id]; !ok {
		return nil, errors.New("not found")
	}
	m.Data[id] = e
	return e, nil
}

func (m *MockCrudDAO) Delete(id string) (bool, error) {
	if _, ok := m.Data[id]; !ok {
		return false, errors.New("not found")
	}
	delete(m.Data, id)
	return true, nil
}

// Para cumplir la interfaz CrudDAO[T]
var _ dao.CrudDAO[MockEntityHandler] = (*MockCrudDAO)(nil)

// ---------- TESTS ----------

func TestCreateEntity(t *testing.T) {
	t.Log("🔍 TEST: Validates POST /test creates a new entity successfully")

	e := echo.New()

	mockDao := NewMockCrudDAO()
	svc := service.NewCrudService(mockDao)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest("POST", "/test",
		strings.NewReader(`{"id":"1","name":"phone"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateEntity(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":"1"`)

	t.Log("✅ Entity successfully created and returned with status 201")
}

func TestCreateEntity_InvalidJSON(t *testing.T) {
	t.Log("🔍 TEST: Ensures invalid JSON returns a 400 Bad Request")

	e := echo.New()

	mockDao := NewMockCrudDAO()
	svc := service.NewCrudService(mockDao)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest("POST", "/test",
		strings.NewReader(`{invalid json`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler.CreateEntity(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "JSON inválido")

	t.Log("✅ Invalid JSON correctly returned 400 Bad Request")
}

func TestGetAllEntities(t *testing.T) {
	t.Log("🔍 TEST: Ensures GET /test returns all stored entities")

	e := echo.New()

	mockDao := NewMockCrudDAO()
	mockDao.Data["a"] = &MockEntityHandler{ID: "a", Name: "Phone"}

	svc := service.NewCrudService(mockDao)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler.GetAllEntities(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"Phone"`)

	t.Log("✅ All entities successfully returned")
}

func TestGetEntityByID(t *testing.T) {
	t.Log("🔍 TEST: Validates GET /test/:id returns the correct entity")

	e := echo.New()

	mockDao := NewMockCrudDAO()
	mockDao.Data["1"] = &MockEntityHandler{ID: "1", Name: "TV"}

	svc := service.NewCrudService(mockDao)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest("GET", "/test/1", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler.GetEntityByID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"TV"`)

	t.Log("✅ Entity successfully fetched by ID")
}

func TestGetEntityByID_NotFound(t *testing.T) {
	t.Log("🔍 TEST: Ensures 404 is returned when requested ID does not exist")

	e := echo.New()

	mockDao := NewMockCrudDAO()
	svc := service.NewCrudService(mockDao)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest("GET", "/test/abc", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	handler.GetEntityByID(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	t.Log("✅ Missing entity correctly returned 404")
}

func TestUpdateEntity(t *testing.T) {
	t.Log("🔍 TEST: Validates PATCH /test/:id updates an existing entity")

	e := echo.New()

	mockDao := NewMockCrudDAO()
	mockDao.Data["1"] = &MockEntityHandler{ID: "1", Name: "Old"}

	svc := service.NewCrudService(mockDao)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest("PATCH", "/test/1",
		strings.NewReader(`{"id":"1","name":"New"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler.UpdateEntity(c)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Contains(t, rec.Body.String(), `"New"`)

	t.Log("✅ Entity updated successfully with status 202")
}

func TestDeleteEntity(t *testing.T) {
	t.Log("🔍 TEST: Validates DELETE /test/:id removes entity successfully")

	e := echo.New()

	mockDao := NewMockCrudDAO()
	mockDao.Data["200"] = &MockEntityHandler{ID: "200", Name: "To Delete"}

	svc := service.NewCrudService(mockDao)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest("DELETE", "/test/200", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("200")

	handler.DeleteEntity(c)

	assert.Equal(t, http.StatusNoContent, rec.Code)

	t.Log("✅ Entity deleted successfully with status 204")
}
