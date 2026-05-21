package repository

import (
	"github.com/rwndy/bookmark-api/internal/domain"
	"gorm.io/gorm"
)

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) domain.BookmarkRepository {
	return &bookmarkRepository{db: db}
}

func (r *bookmarkRepository) Create(bookmark *domain.Bookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *bookmarkRepository) FindByUserID(userID uint) ([]domain.Bookmark, error) {
	var bookmarks []domain.Bookmark
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookmarks).Error
	return bookmarks, err
}

func (r *bookmarkRepository) FindByID(id uint) (*domain.Bookmark, error) {
	var bookmark domain.Bookmark
	err := r.db.First(&bookmark, id).Error
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

func (r *bookmarkRepository) Update(bookmark *domain.Bookmark) error {
	return r.db.Save(bookmark).Error
}

func (r *bookmarkRepository) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Bookmark{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}