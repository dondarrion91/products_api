package utils

import models "project/pkg"

type ProductPayload struct {
	ID              string                       `json:"id"`
	Name            string                       `json:"name" validate:"required"`
	Rate            int                          `json:"rate"`
	Price           float64                      `json:"price" validate:"required,gt=0"`
	Discount        float64                      `json:"discount" validate:"required,gt=0"`
	Installments    int                          `json:"installments" validate:"required,gt=0"`
	Stock           *int                         `json:"stock" validate:"required,gte=0"`
	Details         []models.ProductDetail       `json:"details" validate:"required,dive"`
	Images          []models.Image               `json:"images"`
	SalesNumber     int                          `json:"sales_number"`
	Description     string                       `json:"description"`
	CategoryId      string                       `json:"categoryId" validate:"required"`
	SellerId        string                       `json:"sellerId" validate:"required"`
	Characteristics models.ProductCharacteristic `json:"characteristic"`
}

type ChangeAttributePayload struct {
	ID string `json:"id"`
}
