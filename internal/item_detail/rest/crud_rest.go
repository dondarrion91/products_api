package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"project/internal/item_detail/service"
	"project/internal/item_detail/utils"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	validate = validator.New()
)

type CrudHandler[T any] struct {
	service *service.CrudService[T]
	lock    sync.RWMutex
	clock   sync.Mutex
}

func NewCrudHandler[T any](s *service.CrudService[T]) *CrudHandler[T] {
	return &CrudHandler[T]{service: s}
}

func BindJSON(c echo.Context, entity interface{}) error {
	err := c.Bind(entity)
	if err == nil {
		return nil
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return fmt.Errorf(
			"invalid field '%s': must be %s",
			typeErr.Field,
			typeErr.Type.String(),
		)
	}

	return fmt.Errorf("JSON inválido")
}

func (h *CrudHandler[T]) CreateEntity(c echo.Context) error {
	// Evito problemas de concurrencia
	h.clock.Lock()
	defer h.clock.Unlock()

	var entity T

	if err := BindJSON(c, &entity); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	if err := validate.Struct(entity); err != nil {
		return utils.ValidateBody(c, err)
	}

	createdEntity, err := h.service.RegisterEntity(&entity)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, createdEntity)
}

func (h *CrudHandler[T]) GetAllEntities(c echo.Context) error {
	h.lock.RLock()
	defer h.lock.RUnlock()

	q := c.QueryParam("q")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	entities, err := h.service.FetchEntities(q, limit, offset)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entities)
}

func (h *CrudHandler[T]) GetEntityByID(c echo.Context) error {
	h.lock.RLock()
	defer h.lock.RUnlock()

	id := c.Param("id")

	entity, err := h.service.FetchEntity(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity)
}

func (h *CrudHandler[T]) UpdateEntity(c echo.Context) error {
	// Evito problemas de concurrencia
	h.clock.Lock()
	defer h.clock.Unlock()

	var entity T

	if err := BindJSON(c, &entity); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	id := c.Param("id")

	updatedEntity, err := h.service.PatchEntity(&entity, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusAccepted, updatedEntity)
}

func (h *CrudHandler[T]) DeleteEntity(c echo.Context) error {
	h.clock.Lock()
	defer h.clock.Unlock()

	id := c.Param("id")

	deleted, err := h.service.DeleteEntity(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, deleted)
}
