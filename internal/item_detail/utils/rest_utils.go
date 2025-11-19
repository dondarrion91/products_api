package utils

import (
	"net/http"
	"project/internal/item_detail/repo/datasource/dao"
	"project/internal/item_detail/service"
	models "project/pkg"

	"github.com/labstack/echo/v4"
)

func BuildNotFoundResponse(sellerError error, categoryError error) string {
	if sellerError != nil && categoryError != nil {
		return "Seller and Category doesn't exists"
	}

	if sellerError != nil {
		return "Seller doesn't exists"
	}

	if categoryError != nil {
		return "Category doesn't exists"
	}

	return ""
}

func ChangeAttribute[T any](
	c echo.Context,
	entity ChangeAttributePayload,
	productService *service.ProductService,
	attributeService dao.CrudDAO[T],
	setter func(product *models.Product, ID string),
) error {
	// Validaciones del body binding
	if err := c.Bind(&entity); err != nil {
		return ValidateBody(c, err)
	}

	id := c.Param("id")

	product, err := productService.GetProduct(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	setter(product, entity.ID)

	updatedEntity, err := productService.UpdateProduct(id, product)

	return c.JSON(http.StatusCreated, updatedEntity)
}
