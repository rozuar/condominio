package models

import "time"

type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeVideo FileType = "video"
)

type Galeria struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	EventDate     *time.Time `json:"event_date,omitempty"`
	IsPublic      bool       `json:"is_public"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	ItemsCount    int        `json:"items_count"`
	CreatedBy     *string    `json:"created_by,omitempty"`
	CreatorName   string     `json:"creator_name,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type GaleriaItem struct {
	ID           string    `json:"id"`
	GaleriaID    string    `json:"galeria_id"`
	FileURL      string    `json:"file_url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	FileType     FileType  `json:"file_type"`
	Caption      string    `json:"caption,omitempty"`
	OrderIndex   int       `json:"order_index"`
	CreatedAt    time.Time `json:"created_at"`
}

type GaleriaWithItems struct {
	Galeria
	Items []GaleriaItem `json:"items"`
}

type CreateGaleriaRequest struct {
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	EventDate     *time.Time `json:"event_date,omitempty"`
	IsPublic      bool       `json:"is_public"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
}

type UpdateGaleriaRequest struct {
	Title         *string    `json:"title,omitempty"`
	Description   *string    `json:"description,omitempty"`
	EventDate     *time.Time `json:"event_date,omitempty"`
	IsPublic      *bool      `json:"is_public,omitempty"`
	CoverImageURL *string    `json:"cover_image_url,omitempty"`
}

type AddGaleriaItemRequest struct {
	FileURL      string   `json:"file_url"`
	ThumbnailURL string   `json:"thumbnail_url,omitempty"`
	FileType     FileType `json:"file_type"`
	Caption      string   `json:"caption,omitempty"`
	OrderIndex   int      `json:"order_index"`
}

type UpdateGaleriaItemRequest struct {
	FileURL      *string   `json:"file_url,omitempty"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	FileType     *FileType `json:"file_type,omitempty"`
	Caption      *string   `json:"caption,omitempty"`
	OrderIndex   *int      `json:"order_index,omitempty"`
}

type GaleriaListResponse struct {
	Galerias []Galeria `json:"galerias"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PerPage  int       `json:"per_page"`
}

type GaleriaFilter struct {
	IsPublic *bool
	Page     int
	PerPage  int
}

func (ft FileType) IsValid() bool {
	switch ft {
	case FileTypeImage, FileTypeVideo:
		return true
	}
	return false
}
