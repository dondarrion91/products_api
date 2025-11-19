package main_test

import (
	"fmt"
	models "project/pkg"
	"testing"

	"github.com/google/uuid"
)

// ----------------------------
// Test Init()
// ----------------------------
func TestProduct_Init_GeneratesIDIfEmpty(t *testing.T) {
	t.Log("🔍 TEST: Ensures Init() auto-generates a valid UUID when ID is empty")

	p := &models.Product{
		Name:         "Test",
		Price:        100,
		Discount:     10,
		Installments: 2,
	}

	p.Init()

	if p.ID == "" {
		t.Fatalf("expected ID to be generated, got empty string")
	}

	_, err := uuid.Parse(p.ID)
	if err != nil {
		t.Fatalf("expected valid UUID, got %s", p.ID)
	}

	t.Log("✅ Init() successfully generated UUID")
}

func TestProduct_Init_DoesNotOverrideID(t *testing.T) {
	t.Log("🔍 TEST: Ensures Init() does NOT override an existing ID")

	id := uuid.New().String()
	p := &models.Product{
		ID:           id,
		Name:         "Test",
		Price:        100,
		Discount:     10,
		Installments: 2,
	}

	p.Init()

	if p.ID != id {
		t.Fatalf("expected ID %s, got %s", id, p.ID)
	}

	t.Log("✅ Init() respected the existing ID")
}

func TestProduct_Init_ComputesDerivedPrices(t *testing.T) {
	t.Log("🔍 TEST: Ensures Init() computes InstallmentPrice and DiscountPrice")

	p := &models.Product{
		Name:         "Test",
		Price:        200,
		Discount:     25, // -25%
		Installments: 4,
	}

	p.Init()

	expectedInstallment := 50.0
	expectedDiscount := 150.0 // 200 - 25%

	if p.InstallmentPrice != expectedInstallment {
		t.Fatalf("expected InstallmentPrice %.2f, got %.2f", expectedInstallment, p.InstallmentPrice)
	}

	if p.DiscountPrice != expectedDiscount {
		t.Fatalf("expected DiscountPrice %.2f, got %.2f", expectedDiscount, p.DiscountPrice)
	}

	t.Log("✅ Derived prices computed successfully")
}

// ----------------------------
// Test CalculateInstallmentPrice()
// ----------------------------
func TestCalculateInstallmentPrice(t *testing.T) {
	t.Log("🔍 TEST: Validates CalculateInstallmentPrice() divides price by installments")

	p := &models.Product{Price: 120, Installments: 3}
	value := p.CalculateInstallmentPrice()

	expected := 40.0
	if value != expected {
		t.Fatalf("expected %.2f, got %.2f", expected, value)
	}

	t.Log("✅ Installment price calculated correctly")
}

func TestCalculateInstallmentPrice_NoInstallments(t *testing.T) {
	t.Log("🔍 TEST: Ensures 0 installments returns full price as single payment")

	p := &models.Product{Price: 120, Installments: 0}
	value := p.CalculateInstallmentPrice()

	expected := 120.0
	if value != expected {
		t.Fatalf("expected %.2f, got %.2f", expected, value)
	}

	t.Log("✅ Single payment behavior verified")
}

// ----------------------------
// Test CalculatePriceWithDiscount()
// ----------------------------
func TestCalculatePriceWithDiscount(t *testing.T) {
	t.Log("🔍 TEST: Validates price with discount is correctly computed")

	p := &models.Product{Price: 100, Discount: 20}
	value := p.CalculatePriceWithDiscount()

	expected := 80.0
	if value != expected {
		t.Fatalf("expected %.2f, got %.2f", expected, value)
	}

	t.Log("✅ Discount applied correctly")
}

func TestCalculatePriceWithDiscount_ZeroDiscount(t *testing.T) {
	t.Log("🔍 TEST: Ensures 0% discount returns original price")

	p := &models.Product{Price: 100, Discount: 0}
	value := p.CalculatePriceWithDiscount()

	expected := 100.0
	if value != expected {
		t.Fatalf("expected %.2f, got %.2f", expected, value)
	}

	t.Log("✅ Zero-discount behavior verified")
}

// ----------------------------
// Test ToString()
// ----------------------------
func TestProduct_ToString(t *testing.T) {
	t.Log("🔍 TEST: Ensures ToString() prints formatted product information")

	p := &models.Product{Name: "iPhone", Price: 1500}
	value := p.ToString()

	expected := fmt.Sprintf("Product: iPhone (1500.00)")
	if value != expected {
		t.Fatalf("expected %s, got %s", expected, value)
	}

	t.Log("✅ ToString() formatted output correctly")
}
