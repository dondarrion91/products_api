package rest

import (
	"net/http"
	"project/internal/item_detail/service"
	"project/internal/item_detail/utils"
	models "project/pkg"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

func (h *ProductHandler) GetProduct(c echo.Context) (*models.Product, error) {
	// Evito problemas de concurrencia
	lock.Lock()
	defer lock.Unlock()

	id := c.Param("id")

	product, err := h.service.GetProduct(id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (h *ProductHandler) ChangeCategories(c echo.Context) error {
	// Evito problemas de concurrencia
	lock.Lock()
	defer lock.Unlock()

	var entity utils.ChangeAttributePayload

	categoryService := h.service.GetCategoryService()

	return utils.ChangeAttribute(
		c,
		entity,
		h.service,
		categoryService,
		func(p *models.Product, ID string) {
			p.CategoryId = ID
		},
	)
}

func (h *ProductHandler) AddImages(c echo.Context) error {
	// Evito problemas de concurrencia
	lock.Lock()
	defer lock.Unlock()

	var entity utils.ChangeAttributePayload

	ImageService := h.service.GetImageService()

	return utils.ChangeAttribute(
		c,
		entity,
		h.service,
		ImageService,
		func(p *models.Product, ID string) {
			images := append(p.Images, ID)

			p.Images = images
		},
	)
}

func (h *ProductHandler) ChangeSellers(c echo.Context) error {
	// Evito problemas de concurrencia
	lock.Lock()
	defer lock.Unlock()

	var entity utils.ChangeAttributePayload

	sellerService := h.service.GetSellerService()

	return utils.ChangeAttribute(
		c,
		entity,
		h.service,
		sellerService,
		func(p *models.Product, ID string) {
			p.SellerId = ID
		},
	)
}

func (h *ProductHandler) GetCategories(c echo.Context) error {
	product, err := h.GetProduct(c)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	category, categoryErr := h.service.GetCategoryService().GetByID(product.CategoryId)

	if categoryErr != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": categoryErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, category)
}

func (h *ProductHandler) GetSellers(c echo.Context) error {
	product, err := h.GetProduct(c)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	seller, sellerErr := h.service.GetSellerService().GetByID(product.SellerId)

	if sellerErr != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": sellerErr.Error(),
		})
	}

	return c.JSON(http.StatusOK, seller)
}

func (h *ProductHandler) GetImages(c echo.Context) error {
	product, err := h.GetProduct(c)

	var response []*models.Image

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	for _, id := range product.Images {
		image, imageErr := h.service.GetImageService().GetByID(id)

		if imageErr != nil {
			continue
		}

		response = append(response, image)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetCharacteristic(c echo.Context) error {
	product, err := h.GetProduct(c)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, product.Characteristics)
}

func (h *ProductHandler) GetDetails(c echo.Context) error {
	product, err := h.GetProduct(c)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, product.Details)
}
