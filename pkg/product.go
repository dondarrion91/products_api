package models

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Product struct {
	ID           string          `json:"id"`
	Name         string          `json:"name" validate:"required"`
	Rate         int             `json:"rate"`
	Price        float64         `json:"price" validate:"required,gt=0"`
	Discount     float64         `json:"discount" validate:"required,gt=0"`
	Installments int             `json:"installments" validate:"required,gt=0"`
	Stock        *int            `json:"stock" validate:"required,gte=0"`
	Details      []ProductDetail `json:"details" validate:"required,min=1,dive"`
	Images       []string        `json:"images"`
	SalesNumber  int             `json:"sales_number"`
	Description  string          `json:"description"`

	CategoryId string `json:"categoryId" validate:"required"`
	SellerId   string `json:"sellerId" validate:"required"`

	Characteristics  ProductCharacteristic `json:"characteristics"`
	InstallmentPrice float64               `json:"installmentPrice"`
	DiscountPrice    float64               `json:"discountPrice"`

	// HATEOAS
	Category   HATEOASLink   `json:"category"`
	Seller     HATEOASLink   `json:"seller"`
	ImageLinks []HATEOASLink `json:"image_links"`
}

type HATEOASLink struct {
	Href string `json:"href"`
}

func (p *Product) Init() {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	// Construcción de HATEOAS
	base := os.Getenv("BASE_URL")

	if p.CategoryId != "" {
		p.Category = HATEOASLink{
			Href: fmt.Sprintf("%s/categories/%s", base, p.CategoryId),
		}
	}

	if p.SellerId != "" {
		p.Seller = HATEOASLink{
			Href: fmt.Sprintf("%s/sellers/%s", base, p.SellerId),
		}
	}

	p.ImageLinks = make([]HATEOASLink, 0, len(p.Images))
	for _, imageId := range p.Images {
		p.ImageLinks = append(p.ImageLinks, HATEOASLink{
			Href: fmt.Sprintf("%s/images/%s", base, imageId),
		})
	}

	p.InstallmentPrice = p.CalculateInstallmentPrice()
	p.DiscountPrice = p.CalculatePriceWithDiscount()
}

type ProductCharacteristic struct {
	Name    string          `json:"name" validate:"required"`
	Details []ProductDetail `json:"details" validate:"required,dive"`
}

type ProductDetail struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (p *Product) Create() {}

// ToString devuelve una representación del producto como string.
func (p *Product) ToString() string {
	return fmt.Sprintf("Product: %s (%.2f)", p.Name, p.Price)
}

// CalculateInstallmentPrice calcula el precio por cuota.
func (p *Product) CalculateInstallmentPrice() float64 {
	if p.Installments <= 0 {
		return p.Price
	}
	return p.Price / float64(p.Installments)
}

// CalculatePriceWithDiscount calcula el precio con descuento aplicado.
func (p *Product) CalculatePriceWithDiscount() float64 {
	return p.Price - (p.Price * (p.Discount / 100))
}
