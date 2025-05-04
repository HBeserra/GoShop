package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Media struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"primaryKey;index:idx_media"`
	Namespace string    `json:"namespace" gorm:"index:idx_media"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Filename string `json:"filename"`
	// Url is the public link for the media
	Url string `json:"url"`
	// Type indicates the media's content type, such as image, video, or document.
	Type string `json:"type" gorm:"index:idx_media"`
	// Size represents the size of the media in bytes.
	Size int `json:"size"`
}
