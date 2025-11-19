package main_test

import (
	"path/filepath"
	"testing"

	"project/internal/item_detail/repo/datasource/dal"
	models "project/pkg"

	"github.com/stretchr/testify/assert"
)

func newTempCategoryDAL(t *testing.T) *dal.CrudDAL[models.Category] {
	t.Log("🧪 Setting up a temporary DAL instance for Category tests")

	dir := t.TempDir()
	filename := filepath.Join(dir, "categories")

	return &dal.CrudDAL[models.Category]{
		Filename: filename,
	}
}

func TestCategoryDAL_Create(t *testing.T) {
	t.Log("🔍 TEST: Ensures Create() stores a Category in the JSON file")

	repo := newTempCategoryDAL(t)

	cat := &models.Category{
		ID:   "1",
		Name: "Electrónica",
	}

	created, err := repo.Create(cat)
	assert.NoError(t, err)
	assert.Equal(t, "1", created.ID)

	all, err := repo.GetAll("", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, all, 1)
	assert.Equal(t, "Electrónica", all[0].Name)

	t.Log("✅ Category successfully created and persisted")
}

func TestCategoryDAL_GetByID(t *testing.T) {
	t.Log("🔍 TEST: Validates GetByID() returns an existing Category")

	repo := newTempCategoryDAL(t)

	repo.Create(&models.Category{ID: "10", Name: "Ropa"})

	cat, err := repo.GetByID("10")
	assert.NoError(t, err)
	assert.Equal(t, "Ropa", cat.Name)

	t.Log("✅ Category retrieved correctly by ID")
}

func TestCategoryDAL_GetByID_NotFound(t *testing.T) {
	t.Log("🔍 TEST: Ensures GetByID() returns an error when the ID does not exist")

	repo := newTempCategoryDAL(t)

	_, err := repo.GetByID("no-existe")
	assert.Error(t, err)

	t.Log("✅ Correctly returned error for unknown Category ID")
}

func TestCategoryDAL_Update(t *testing.T) {
	t.Log("🔍 TEST: Ensures Update() modifies the stored Category correctly")

	repo := newTempCategoryDAL(t)

	repo.Create(&models.Category{ID: "2", Name: "Deportes"})

	updated, err := repo.Update(&models.Category{
		ID:   "2",
		Name: "Camping",
	}, "2")

	assert.NoError(t, err)
	assert.Equal(t, "Camping", updated.Name)
	all, _ := repo.GetAll("", 10, 0)
	assert.Equal(t, "Camping", all[0].Name)

	t.Log("✅ Category updated successfully")
}

func TestCategoryDAL_Delete(t *testing.T) {
	t.Log("🔍 TEST: Ensures Delete() removes a Category when the ID exists")

	repo := newTempCategoryDAL(t)

	repo.Create(&models.Category{ID: "7", Name: "Autos"})

	_, err := repo.Delete("7")
	assert.NoError(t, err)

	all, _ := repo.GetAll("", 10, 0)
	assert.Len(t, all, 0)

	t.Log("✅ Category deleted successfully")
}

func TestCategoryDAL_Delete_NotFound(t *testing.T) {
	t.Log("🔍 TEST: Ensures Delete() returns an error when the ID does not exist")

	repo := newTempCategoryDAL(t)

	_, err := repo.Delete("nope")
	assert.Error(t, err)

	t.Log("✅ Correctly returned error for deleting non-existing Category")
}
