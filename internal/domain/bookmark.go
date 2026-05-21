package domain

import "time"

type Bookmark struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	URL       string    `json:"url" gorm:"not null"`
	Title     string    `json:"title"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookmarkRepository interface {
	Create(bookmark *Bookmark) error
	FindByUserID(userID uint) ([]Bookmark, error)
	FindByID(id uint) (*Bookmark, error)
	Update(bookmark *Bookmark) error
	Delete(id uint, userID uint) error
}

type BookmarkService interface {
	CreateBookmark(userID uint, url, title, tags string) (*Bookmark, error)
	GetUserBookmarks(userID uint) ([]Bookmark, error)
	UpdateBookmark(id, userID uint, url, title, tags string) (*Bookmark, error)
	DeleteBookmark(id, userID uint) error
}