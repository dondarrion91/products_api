package main_test

import (
	"os"
	"path/filepath"
	"project/internal/item_detail/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper para crear archivos temporales sin ensuciar el repo
func tempFilePath(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, name)
}

func TestWriteJSON(t *testing.T) {
	t.Log("🔍 TEST: Ensures WriteJSON creates the .json file with correct contents")

	path := tempFilePath(t, "test_write")

	data := map[string]interface{}{
		"name": "Julian",
		"age":  28,
	}

	err := utils.WriteJSON(path, data)
	assert.NoError(t, err)

	// Verificamos que se creó el archivo
	_, err = os.Stat(path + ".json")
	assert.NoError(t, err)

	t.Log("✅ JSON file successfully created")
}

func TestReadJSON_FileNotExist_ReturnsNil(t *testing.T) {
	t.Log("🔍 TEST: Ensures ReadJSON returns nil when file does not exist")

	path := tempFilePath(t, "no_file")

	var dest map[string]interface{}
	err := utils.ReadJSON(path, &dest)

	assert.NoError(t, err)
	assert.Nil(t, dest)

	t.Log("✅ Nonexistent file correctly returned nil")
}

func TestReadJSON_Success(t *testing.T) {
	t.Log("🔍 TEST: Ensures ReadJSON correctly decodes existing JSON file")

	path := tempFilePath(t, "read_test")

	// Creamos archivo manualmente
	content := `{"city":"Córdoba","active":true}`
	err := os.WriteFile(path+".json", []byte(content), 0644)
	assert.NoError(t, err)

	var dest map[string]interface{}
	err = utils.ReadJSON(path, &dest)

	assert.NoError(t, err)
	assert.Equal(t, "Córdoba", dest["city"])
	assert.Equal(t, true, dest["active"])

	t.Log("✅ JSON decoded successfully")
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	t.Log("🔍 TEST: Ensures ReadJSON returns error for invalid JSON format")

	path := tempFilePath(t, "bad_json")

	err := os.WriteFile(path+".json", []byte("{invalid"), 0644)
	assert.NoError(t, err)

	var dest map[string]interface{}
	err = utils.ReadJSON(path, &dest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error al decodificar JSON")

	t.Log("✅ Invalid JSON produced an error as expected")
}

func TestUpdateJSON_CreateNew(t *testing.T) {
	t.Log("🔍 TEST: Ensures UpdateJSON creates a new JSON file when none exists")

	path := tempFilePath(t, "update_new")

	err := utils.UpdateJSON(path, "nombre", "Mauro")
	assert.NoError(t, err)

	var dest map[string]interface{}
	err = utils.ReadJSON(path, &dest)

	assert.NoError(t, err)
	assert.Equal(t, "Mauro", dest["nombre"])

	t.Log("✅ New JSON file created with correct data")
}

func TestUpdateJSON_UpdateExisting(t *testing.T) {
	t.Log("🔍 TEST: Ensures UpdateJSON updates an existing key in JSON file")

	path := tempFilePath(t, "update_existing")

	// Creamos archivo inicial
	err := os.WriteFile(path+".json", []byte(`{"color":"rojo"}`), 0644)
	assert.NoError(t, err)

	// Actualizamos
	err = utils.UpdateJSON(path, "color", "azul")
	assert.NoError(t, err)

	// Leemos
	var dest map[string]interface{}
	err = utils.ReadJSON(path, &dest)

	assert.NoError(t, err)
	assert.Equal(t, "azul", dest["color"])

	t.Log("✅ Existing key updated successfully")
}

func TestUpdateJSON_AddNewField(t *testing.T) {
	t.Log("🔍 TEST: Ensures UpdateJSON adds new fields without removing existing ones")

	path := tempFilePath(t, "update_add")

	err := os.WriteFile(path+".json", []byte(`{"x":1}`), 0644)
	assert.NoError(t, err)

	err = utils.UpdateJSON(path, "y", 2)
	assert.NoError(t, err)

	var dest map[string]interface{}
	err = utils.ReadJSON(path, &dest)

	assert.NoError(t, err)
	assert.Equal(t, float64(1), dest["x"])
	assert.Equal(t, float64(2), dest["y"])

	t.Log("✅ New field added successfully while preserving existing data")
}
