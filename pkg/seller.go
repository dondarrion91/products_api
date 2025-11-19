package models

import "github.com/google/uuid"

type Seller struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Verified bool   `json:"verified"`
}

func (p *Seller) Init() {
	if p.ID == "" {
		p.ID = uuid.New().String()
		p.Verified = false
	}
}
