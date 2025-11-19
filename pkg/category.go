package models

import "github.com/google/uuid"

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required"`
}

func (p *Category) Init() {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
}
