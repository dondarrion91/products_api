package dal

import (
	"fmt"
	"project/internal/item_detail/utils"
	"reflect"
	"strings"
)

// CrudDAL es una implementación genérica CRUD basada en archivos JSON.
// El archivo almacena un array de entidades del tipo T.
type CrudDAL[T any] struct {
	Filename string // Ruta al archivo JSON donde se guardan las entidades
}

type Initializable interface {
	Init()
}

// Create agrega una nueva entidad al archivo JSON.
func (u *CrudDAL[T]) Create(entity *T) (*T, error) {
	var data []T

	if v, ok := any(entity).(Initializable); ok {
		v.Init()
	}

	if err := utils.ReadJSON(u.Filename, &data); err != nil {
		data = []T{}
	}

	data = append(data, *entity)

	if err := utils.WriteJSON(u.Filename, data); err != nil {
		return nil, fmt.Errorf("Can't save the entity with error: %w", err)
	}

	return entity, nil
}

// GetByID busca una entidad con el campo `ID` igual al solicitado.
func (u *CrudDAL[T]) GetByID(uid string) (*T, error) {
	var data []T

	if err := utils.ReadJSON(u.Filename, &data); err != nil {
		return nil, fmt.Errorf("error writing JSON: %w", err)
	}

	for _, item := range data {
		v := reflect.ValueOf(item)

		if v, ok := any(&item).(Initializable); ok {
			v.Init()
		}

		idField := v.FieldByName("ID")

		if idField.IsValid() && idField.Kind() == reflect.String && idField.String() == uid {
			return &item, nil
		}
	}

	return nil, fmt.Errorf("Can't find entity with UID %s", uid)
}

// GetAll devuelve todas las entidades del JSON.
func (u *CrudDAL[T]) GetAll(q string, limit int, offset int) ([]*T, error) {
	var data []T

	if err := utils.ReadJSON(u.Filename, &data); err != nil {
		return nil, fmt.Errorf("error writing JSON: %w", err)
	}

	// Default pagination values
	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 10
	}

	filtered := []*T{}

	// === 1) Apply filtering ===
	for i := range data {
		item := &data[i]

		v := reflect.ValueOf(item).Elem()
		nameField := v.FieldByName("Name")

		// Apply substring filter "q"
		if q != "" {
			if !nameField.IsValid() || nameField.Kind() != reflect.String {
				continue
			}

			name := strings.ToLower(nameField.String())
			search := strings.ToLower(q)

			if !strings.Contains(name, search) {
				continue
			}
		}

		// Initialize entity if needed
		if vi, ok := any(item).(Initializable); ok {
			vi.Init()
		}

		filtered = append(filtered, item)
	}

	// === 2) Apply pagination ===
	total := len(filtered)

	if offset >= total {
		// Nothing to return
		return []*T{}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	paginated := filtered[offset:end]

	return paginated, nil
}

// Update reemplaza los campos no vacíos de una entidad existente (por ID).
func (u *CrudDAL[T]) Update(entity *T, id string) (*T, error) {
	var data []T

	if err := utils.ReadJSON(u.Filename, &data); err != nil {
		return nil, fmt.Errorf("error writing JSON: %w", err)
	}

	for i := range data {
		v := reflect.ValueOf(data[i])
		idField := v.FieldByName("ID")
		if idField.IsValid() && idField.Kind() == reflect.String && idField.String() == id {
			updated, wasUpdated := updateData(entity, &data[i])

			if !wasUpdated {
				return nil, fmt.Errorf("Update failed: invalid parameters or no parameters provided.")
			}

			if err := utils.WriteJSON(u.Filename, data); err != nil {
				return nil, fmt.Errorf("error writing JSON: %w", err)
			}

			if v, ok := any(updated).(Initializable); ok {
				v.Init()
			}

			return updated, nil
		}
	}

	return nil, fmt.Errorf("Can't find entity with ID %v", id)
}

// Delete elimina una entidad por ID.
func (u *CrudDAL[T]) Delete(id string) (bool, error) {
	var data []T

	if err := utils.ReadJSON(u.Filename, &data); err != nil {
		return false, fmt.Errorf("error reading JSON: %w", err)
	}

	found := false
	newData := make([]T, 0)

	for _, item := range data {
		v := reflect.ValueOf(item)
		idField := v.FieldByName("ID")
		if idField.IsValid() && idField.Kind() == reflect.String && idField.String() == id {
			found = true
			continue // salta el eliminado
		}
		newData = append(newData, item)
	}

	if !found {
		return false, fmt.Errorf("Can't find entity with ID %v", id)
	}

	if err := utils.WriteJSON(u.Filename, newData); err != nil {
		return false, fmt.Errorf("error writing JSON: %w", err)
	}

	return true, nil
}

// updateData reemplaza los campos no vacíos o no cero del origen en el destino.
func updateData[T any](entity *T, existingEntity *T) (*T, bool) {
	valSrc := reflect.ValueOf(entity).Elem()
	valDst := reflect.ValueOf(existingEntity).Elem()

	wasUpdated := false

	for i := 0; i < valSrc.NumField(); i++ {

		fieldName := valSrc.Type().Field(i).Name

		if fieldName == "ID" {
			continue
		}

		srcField := valSrc.Field(i)
		dstField := valDst.Field(i)

		if srcField.IsValid() && !srcField.IsZero() && dstField.CanSet() {
			dstField.Set(srcField)
			wasUpdated = true
		}
	}

	return existingEntity, wasUpdated
}
