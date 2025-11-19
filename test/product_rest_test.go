package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"project/internal/item_detail/repo/datasource/dal"
	"project/internal/item_detail/rest"
	"project/internal/item_detail/service"
	"project/internal/item_detail/utils"
	models "project/pkg"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

/* ===========================================================
   Helper: crea producto válido usado en todos los tests
   =========================================================== */

func createTestProduct() models.Product {
	stock := 10
	return models.Product{
		ID:           "1",
		Name:         "Test",
		Price:        100,
		Discount:     10,
		Installments: 1,
		Images: []string{
			"9afa306b-1a10-472f-965b-09dd511d56d1",
		},
		Stock: &stock,
		Details: []models.ProductDetail{
			{Name: "Size", Description: "Large"},
		},
		CategoryId: "100",
		SellerId:   "200",
		Characteristics: models.ProductCharacteristic{
			Name: "Specs",
			Details: []models.ProductDetail{
				{Name: "CPU", Description: "Fast"},
			},
		},
	}
}

/* ===========================================================
   Helper: reescribe JSON y lo restaura o elimina al final
   =========================================================== */

func SafeRewriteJSON(t *testing.T, filename string, data interface{}) {
	// Ver si existía
	original, err := os.ReadFile(filename)
	existed := (err == nil)

	// WriteJSON requiere nombre SIN .json
	base := strings.TrimSuffix(filename, ".json")

	err = utils.WriteJSON(base, data)
	assert.NoError(t, err)

	// Limpieza automática al finalizar el test
	t.Cleanup(func() {
		if existed {
			// Restaurar archivo original
			_ = ioutil.WriteFile(filename, original, 0644)
		} else {
			// Borrar archivo creado por el test
			_ = os.Remove(filename)
		}
	})
}

/* ===========================================================
   Test GetProduct()
   =========================================================== */

/* ===========================================================
   Test GetProduct()
   =========================================================== */

func TestGetProduct(t *testing.T) {
	t.Log("🔍 TEST: Validates GetProduct returns product data for a valid ID")

	e := echo.New()

	// Crear producto temporal
	product := createTestProduct()
	SafeRewriteJSON(t, "Product.json", []models.Product{product})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	req := httptest.NewRequest("GET", "/products/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("1")

	p, err := handler.GetProduct(c)
	assert.NoError(t, err)
	assert.Equal(t, "1", p.ID)
	assert.Equal(t, "Test", p.Name)

	t.Log("✅ Product successfully retrieved")
}

/* ===========================================================
   Test GetCategories()
   =========================================================== */

func TestGetCategories(t *testing.T) {
	t.Log("🔍 TEST: Ensures categories for a product are correctly returned")

	e := echo.New()

	product := createTestProduct()

	SafeRewriteJSON(t, "Product.json", []models.Product{product})
	SafeRewriteJSON(t, "Category.json", []models.Category{
		{ID: "100", Name: "Electronics"},
	})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	req := httptest.NewRequest("GET", "/products/1/categories", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetCategories(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"Electronics"`)

	t.Log("✅ Categories returned successfully")
}

/* ===========================================================
   Test GetSellers()
   =========================================================== */

func TestGetSellers(t *testing.T) {
	t.Log("🔍 TEST: Ensures sellers related to a product are correctly resolved")

	e := echo.New()

	product := createTestProduct()

	SafeRewriteJSON(t, "Product.json", []models.Product{product})
	SafeRewriteJSON(t, "Seller.json", []models.Seller{
		{ID: "200", Name: "TestSeller"},
	})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	req := httptest.NewRequest("GET", "/products/1/sellers", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetSellers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"TestSeller"`)

	t.Log("✅ Sellers returned successfully")
}

/* ===========================================================
   Test GetImages()
   =========================================================== */

func TestGetImages(t *testing.T) {
	t.Log("🔍 TEST: Ensures images linked to the product are correctly fetched")

	e := echo.New()

	product := createTestProduct()

	SafeRewriteJSON(t, "Product.json", []models.Product{product})
	SafeRewriteJSON(t, "Image.json", []models.Image{
		{ID: "9afa306b-1a10-472f-965b-09dd511d56d1", Name: "TestImages"},
	})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	req := httptest.NewRequest("GET", "/products/1/images", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetImages(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"TestImages"`)

	t.Log("✅ Images returned successfully")
}

/* ===========================================================
   Test GetCharacteristic()
   =========================================================== */

func TestGetCharacteristic(t *testing.T) {
	t.Log("🔍 TEST: Validates product characteristic retrieval")

	e := echo.New()

	product := createTestProduct()
	SafeRewriteJSON(t, "Product.json", []models.Product{product})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	req := httptest.NewRequest("GET", "/products/1/characteristic", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.GetCharacteristic(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"CPU"`)

	t.Log("✅ Characteristic returned successfully")
}

/* ===========================================================
   Test ChangeCategories()
   =========================================================== */

func TestChangeCategories(t *testing.T) {
	t.Log("🔍 TEST: Confirms product category is updated correctly")

	e := echo.New()

	product := createTestProduct()

	SafeRewriteJSON(t, "Product.json", []models.Product{product})
	SafeRewriteJSON(t, "Category.json", []models.Category{
		{ID: "100", Name: "Electronics"},
		{ID: "999", Name: "Games"},
	})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	body := `{"id":"999"}`
	req := httptest.NewRequest("PATCH", "/products/1/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.ChangeCategories(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	updatedBytes, _ := ioutil.ReadFile("Product.json")
	assert.Contains(t, string(updatedBytes), `"categoryId": "999"`)

	t.Log("✅ Category updated successfully")
}

/* ===========================================================
   Test ChangeSellers()
   =========================================================== */

func TestChangeSellers(t *testing.T) {
	t.Log("🔍 TEST: Confirms product seller is updated correctly")

	e := echo.New()

	product := createTestProduct()

	SafeRewriteJSON(t, "Product.json", []models.Product{product})
	SafeRewriteJSON(t, "Seller.json", []models.Seller{
		{ID: "200", Name: "Old"},
		{ID: "333", Name: "NewSeller"},
	})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	body := `{"id":"333"}`
	req := httptest.NewRequest("PATCH", "/products/1/sellers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.ChangeSellers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	updatedBytes, _ := ioutil.ReadFile("Product.json")
	assert.Contains(t, string(updatedBytes), `"sellerId": "333"`)

	t.Log("✅ Seller updated successfully")
}

/* ===========================================================
   Tests for GetAllEntities() with query param ?q=
   =========================================================== */

func TestGetAllProducts_FilterByQuery_SingleMatch(t *testing.T) {
	t.Log("🔍 TEST: Returns only the product that matches the ?q= filter")

	e := echo.New()

	products := []models.Product{
		{ID: "1", Name: "Laptop Gamer"},
		{ID: "2", Name: "Mouse"},
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products?q=gamer", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"Laptop Gamer"`)
	assert.NotContains(t, body, `"Mouse"`)

	t.Log("✅ Correctly returned only the product matching 'gamer'")
}

func TestGetAllProducts_FilterByQuery_MultipleMatch(t *testing.T) {
	t.Log("🔍 TEST: Returns all products matching the substring in ?q=")

	e := echo.New()

	products := []models.Product{
		{ID: "1", Name: "MacBook Pro"},
		{ID: "2", Name: "Mac Mini"},
		{ID: "3", Name: "Monitor Samsung"},
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products?q=mac", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()
	assert.Contains(t, body, `"MacBook Pro"`)
	assert.Contains(t, body, `"Mac Mini"`)
	assert.NotContains(t, body, `"Monitor Samsung"`)

	t.Log("✅ Correctly returned all products matching 'mac'")
}

func TestGetAllProducts_FilterByQuery_NoMatch(t *testing.T) {
	t.Log("🔍 TEST: When no items match, the endpoint must return an empty array []")

	e := echo.New()

	products := []models.Product{
		{ID: "1", Name: "Camera"},
		{ID: "2", Name: "Speaker"},
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products?q=xyz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "[]", strings.TrimSpace(rec.Body.String()))

	t.Log("✅ No matches: returned [] as expected")
}

func TestGetAllProducts_NoQueryParam_ReturnsAll(t *testing.T) {
	t.Log("🔍 TEST: Without ?q=, the endpoint must return all products")

	e := echo.New()

	products := []models.Product{
		{ID: "1", Name: "Table"},
		{ID: "2", Name: "Chair"},
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()
	assert.Contains(t, body, `"Table"`)
	assert.Contains(t, body, `"Chair"`)

	t.Log("✅ Returned all products when ?q= is not provided")
}

func TestGetAllProducts_FilterByQuery_CaseInsensitive(t *testing.T) {
	t.Log("🔍 TEST: Query filter must be case-insensitive")

	e := echo.New()

	products := []models.Product{
		{ID: "1", Name: "Monitor LG"},
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products?q=monitor", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"Monitor LG"`)

	t.Log("✅ Successfully matched ignoring case differences")
}

func TestGetAllProducts_Pagination_DefaultLimit(t *testing.T) {
	t.Log("🔍 TEST: Ensures default pagination returns the first 10 items")

	e := echo.New()

	// Create 15 items
	products := []models.Product{}
	for i := 1; i <= 15; i++ {
		products = append(products, models.Product{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("Item %d", i)})
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products?limit=10&offset=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()

	// Items 1 to 10 must appear
	for i := 1; i <= 10; i++ {
		assert.Contains(t, body, fmt.Sprintf("Item %d", i))
	}

	// Items 11 to 15 must NOT appear
	for i := 11; i <= 15; i++ {
		assert.NotContains(t, body, fmt.Sprintf("Item %d", i))
	}

	t.Log("✅ Returned only first 10 items by default pagination")
}

func TestGetAllProducts_Pagination_OffsetWorks(t *testing.T) {
	t.Log("🔍 TEST: Ensures offset skips items correctly")

	e := echo.New()

	// Create 10 items
	products := []models.Product{}
	for i := 1; i <= 10; i++ {
		products = append(products, models.Product{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("Item %d", i)})
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	// offset=5 → should return items 6–10
	req := httptest.NewRequest(http.MethodGet, "/products?offset=5&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)
	assert.NoError(t, err)

	body := rec.Body.String()

	// Items 1 to 5 should NOT appear
	for i := 1; i <= 5; i++ {
		assert.NotContains(t, body, fmt.Sprintf(`"Item %d"`, i))
	}

	// Items 6 to 10 should appear
	for i := 6; i <= 10; i++ {
		assert.Contains(t, body, fmt.Sprintf("Item %d", i))
	}

	t.Log("✅ Offset correctly skipped the first 5 items")
}

func TestGetAllProducts_Pagination_LimitWorks(t *testing.T) {
	t.Log("🔍 TEST: Ensures custom limit restricts returned items")

	e := echo.New()

	// Create 20 items
	products := []models.Product{}
	for i := 1; i <= 20; i++ {
		products = append(products, models.Product{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("Item %d", i)})
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products?limit=3&offset=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)
	assert.NoError(t, err)

	body := rec.Body.String()

	assert.Contains(t, body, `"Item 1"`)
	assert.Contains(t, body, `"Item 2"`)
	assert.Contains(t, body, `"Item 3"`)

	assert.NotContains(t, body, `"Item 4"`)

	t.Log("✅ Limit successfully restricted the response to 3 items")
}

func TestGetAllProducts_Pagination_OffsetOutOfRange(t *testing.T) {
	t.Log("🔍 TEST: When offset exceeds dataset size, an empty array must be returned")

	e := echo.New()

	products := []models.Product{
		{ID: "1", Name: "Item 1"},
		{ID: "2", Name: "Item 2"},
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/products?offset=10&limit=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)
	assert.NoError(t, err)
	assert.Equal(t, "[]", strings.TrimSpace(rec.Body.String()))

	t.Log("✅ Offset out of range returned empty array []")
}

func TestGetAllProducts_Pagination_DefaultsApplied(t *testing.T) {
	t.Log("🔍 TEST: Negative or missing pagination params should apply defaults (limit=10, offset=0)")

	e := echo.New()

	products := []models.Product{}
	for i := 1; i <= 12; i++ {
		products = append(products, models.Product{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("Item %d", i)})
	}
	SafeRewriteJSON(t, "Product.json", products)

	productDal := dal.NewProductDAL()
	svc := service.NewCrudService(productDal)
	handler := rest.NewCrudHandler(svc)

	// No limit and no offset
	req := httptest.NewRequest(http.MethodGet, "/products?limit=-1&offset=-5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetAllEntities(c)
	assert.NoError(t, err)

	body := rec.Body.String()

	// Must only contain first 10 items
	for i := 1; i <= 10; i++ {
		assert.Contains(t, body, fmt.Sprintf("Item %d", i))
	}

	// Must NOT contain items 11 or 12
	assert.NotContains(t, body, "Item 11")
	assert.NotContains(t, body, "Item 12")

	t.Log("✅ Default pagination (limit=10, offset=0) was applied correctly")
}

func createProductForImages() models.Product {
	stock := 10
	return models.Product{
		ID:           "prod-1",
		Name:         "Camera",
		Price:        100,
		Discount:     5,
		Installments: 1,
		Stock:        &stock,
		Images: []string{
			"existing-img",
		},
		CategoryId: "cat-1",
		SellerId:   "seller-1",
		Details: []models.ProductDetail{
			{Name: "Spec", Description: "Info"},
		},
		Characteristics: models.ProductCharacteristic{
			Name: "Specs",
			Details: []models.ProductDetail{
				{Name: "CPU", Description: "Fast"},
			},
		},
	}
}

/* ===========================================================
   TEST 1 — Add new image to product
   =========================================================== */

func TestAddImages_AddsNewImage(t *testing.T) {
	t.Log("🔍 TEST: Ensures AddImages() appends a new image correctly")

	e := echo.New()

	// --- Product with 1 initial image ---
	product := createProductForImages()
	SafeRewriteJSON(t, "Product.json", []models.Product{product})

	// --- Image exists ---
	SafeRewriteJSON(t, "Image.json", []models.Image{
		{ID: "existing-img", Name: "Existing"},
		{ID: "new-image", Name: "NewImage"},
	})

	// DAL reales
	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	// Body con el ID de la imagen a agregar
	body := `{"id":"new-image"}`

	req := httptest.NewRequest("POST", "/products/prod-1/images", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("prod-1")

	err := handler.AddImages(c)
	assert.NoError(t, err)

	// Validate DB (JSON)
	updatedBytes, _ := ioutil.ReadFile("Product.json")
	bodyStr := string(updatedBytes)

	assert.Contains(t, bodyStr, `"existing-img"`)
	assert.Contains(t, bodyStr, `"new-image"`)

	t.Log("✅ Image appended successfully")
}

/* ===========================================================
   TEST 2 — Product not found
   =========================================================== */

func TestAddImages_ProductNotFound(t *testing.T) {
	t.Log("🔍 TEST: Ensures AddImages() returns 404 for missing product")

	e := echo.New()

	// Product JSON vacío
	SafeRewriteJSON(t, "Product.json", []models.Product{})

	SafeRewriteJSON(t, "Image.json", []models.Image{
		{ID: "img-1", Name: "TestImage"},
	})

	productDal := dal.NewProductDAL()
	sellerDal := dal.NewSellerDAL()
	categoryDal := dal.NewCategoryDAL()
	imageDal := dal.NewImageDAL()

	svc := service.NewProductService(productDal, sellerDal, categoryDal, imageDal)
	handler := rest.NewProductHandler(svc)

	body := `{"id":"img-1"}`

	req := httptest.NewRequest("POST", "/products/prod-1/images", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("prod-1")

	err := handler.AddImages(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	t.Log("✅ Product not found correctly returned 404")
}
