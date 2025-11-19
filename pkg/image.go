package models

import "github.com/google/uuid"

type Image struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required"`
}

func (p *Image) Init() {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
}
