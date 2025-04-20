package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`

	// Url is the public link for the media
	Url string `json:"url"`
	// Type indicates the media's content type, such as image, video, or document.
	Type string `json:"type"`
	// Size represents the size of the media in bytes.
	Size int `json:"size"`
}
