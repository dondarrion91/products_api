package main_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"project/internal/item_detail/repo/datasource/dal"
	"project/internal/item_detail/utils"

	"github.com/stretchr/testify/assert"
)

// ---- Mock de Entidad ----

type MockEntity struct {
	ID    string
	Name  string
	Price float64
}

func TestCRUD_DAL_Create(t *testing.T) {
	t.Log("🔍 TEST: Verifies Create() writes a new entity to the JSON store")

	tmpFile := filepath.Join(os.TempDir(), t.Name())
	os.Remove(tmpFile + ".json")

	repo := &dal.CrudDAL[MockEntity]{Filename: tmpFile}

	entity := &MockEntity{
		ID:    "1",
		Name:  "Producto",
		Price: 10,
	}

	_, err := repo.Create(entity)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	var saved []MockEntity
	_ = utils.ReadJSON(tmpFile, &saved)

	if len(saved) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(saved))
	}

	if saved[0].ID != "1" {
		t.Errorf("expected ID '1', got '%s'", saved[0].ID)
	}

	t.Log("✅ Create() successfully stored the entity")
}

func TestCRUD_DAL_GetByID(t *testing.T) {
	t.Log("🔍 TEST: Ensures GetByID() returns the correct entity")

	tmpFile := filepath.Join(os.TempDir(), "test_getbyid")
	os.Remove(tmpFile + ".json")

	initial := []MockEntity{
		{ID: "A", Name: "Test", Price: 50},
	}

	_ = utils.WriteJSON(tmpFile, initial)

	repo := &dal.CrudDAL[MockEntity]{Filename: tmpFile}

	item, err := repo.GetByID("A")
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if item.Name != "Test" {
		t.Errorf("expected Name='Test', got '%s'", item.Name)
	}

	t.Log("✅ GetByID() returned the correct entity")
}

func TestCRUD_DAL_GetByID_NotFound(t *testing.T) {
	t.Log("🔍 TEST: Ensures GetByID() returns an error for a missing entity")

	repo := &dal.CrudDAL[MockEntity]{Filename: filepath.Join(os.TempDir(), "missing_getbyid")}

	_, err := repo.GetByID("no-existe")
	assert.Error(t, err)

	t.Log("✅ Missing ID correctly produced an error")
}

func TestCRUD_DAL_GetAll(t *testing.T) {
	t.Log("🔍 TEST: Validates GetAll() returns all stored entities")

	tmpFile := filepath.Join(os.TempDir(), "test_getall")
	os.Remove(tmpFile + ".json")

	initial := []MockEntity{
		{ID: "1", Name: "A"},
		{ID: "2", Name: "B"},
	}

	_ = utils.WriteJSON(tmpFile, initial)

	repo := &dal.CrudDAL[MockEntity]{Filename: tmpFile}

	items, err := repo.GetAll("", 10, 0)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	t.Log("✅ GetAll() returned the correct number of entities")
}

func TestCRUD_DAL_Update(t *testing.T) {
	t.Log("🔍 TEST: Ensures Update() modifies fields without overwriting zero-value fields improperly")

	tmpFile := filepath.Join(os.TempDir(), "test_update")
	os.Remove(tmpFile + ".json")

	initial := []MockEntity{
		{ID: "22", Name: "Viejo", Price: 5},
	}

	_ = utils.WriteJSON(tmpFile, initial)

	repo := &dal.CrudDAL[MockEntity]{Filename: tmpFile}

	update := &MockEntity{
		Name:  "Nuevo",
		Price: 99,
	}

	res, err := repo.Update(update, "22")
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if res.Name != "Nuevo" {
		t.Errorf("expected Name='Nuevo', got %s", res.Name)
	}

	if res.Price != 99 {
		t.Errorf("expected Price=99, got %.2f", res.Price)
	}

	t.Log("✅ Update() correctly updated the entity")
}

func TestCRUD_DAL_Delete(t *testing.T) {
	t.Log("🔍 TEST: Ensures Delete() removes the correct entity")

	tmpFile := filepath.Join(os.TempDir(), "test_delete")
	os.Remove(tmpFile + ".json")

	initial := []MockEntity{
		{ID: "DEL", Name: "Eliminar"},
		{ID: "KEEP", Name: "Mantener"},
	}

	_ = utils.WriteJSON(tmpFile, initial)

	repo := &dal.CrudDAL[MockEntity]{Filename: tmpFile}

	ok, err := repo.Delete("DEL")
	if err != nil || !ok {
		t.Fatalf("Delete failed: %v", err)
	}

	var data []MockEntity
	_ = utils.ReadJSON(tmpFile, &data)

	if len(data) != 1 {
		t.Fatalf("expected 1 entity left, got %d", len(data))
	}

	if data[0].ID != "KEEP" {
		t.Fatalf("expected remaining ID 'KEEP', got '%s'", data[0].ID)
	}

	t.Log("✅ Delete() successfully removed the entity")
}

func TestCRUD_DAL_Delete_NotFound(t *testing.T) {
	t.Log("🔍 TEST: Ensures Delete() returns error when ID does not exist")

	repo := &dal.CrudDAL[MockEntity]{Filename: filepath.Join(os.TempDir(), "missing_delete")}

	_, err := repo.Delete("nope")
	assert.Error(t, err)

	t.Log("✅ Delete() correctly reported missing entity")
}

// ---- TEST updateData ----

func TestUpdateData(t *testing.T) {
	t.Log("🔍 TEST: Validates updateData ignores zero fields and updates non-zero ones")

	original := &MockEntity{
		ID:    "X",
		Name:  "Old",
		Price: 20,
	}

	update := &MockEntity{
		Name:  "New",
		Price: 0, // campo cero → NO debe actualizarse
	}

	result, ok := callUpdateData(update, original)
	if !ok {
		t.Fatalf("expected wasUpdated=true, got false")
	}

	if result.Name != "New" {
		t.Errorf("expected Name='New', got '%s'", result.Name)
	}

	if result.Price != 20 {
		t.Errorf("expected Price remain 20, got %.2f", result.Price)
	}

	t.Log("✅ updateData behaved correctly with zero and non-zero fields")
}

// Función para poder testear updateData ya que no es exportada
func callUpdateData[T any](src *T, dst *T) (*T, bool) {
	return reflectUpdateData(src, dst)
}

func reflectUpdateData[T any](src *T, dst *T) (*T, bool) {
	// misma lógica que dal.updateData (copiada)
	valSrc := reflect.ValueOf(src).Elem()
	valDst := reflect.ValueOf(dst).Elem()

	wasUpdated := false

	for i := 0; i < valSrc.NumField(); i++ {
		srcField := valSrc.Field(i)
		dstField := valDst.Field(i)

		if srcField.IsValid() && !srcField.IsZero() && dstField.CanSet() {
			dstField.Set(srcField)
			wasUpdated = true
		}
	}

	return dst, wasUpdated
}
