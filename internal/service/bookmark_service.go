package service

import (
	"net/url"

	"github.com/rwndy/bookmark-api/internal/domain"
)

type bookmarkService struct {
	repo domain.BookmarkRepository
}

func NewBookmarkService(repo domain.BookmarkRepository) domain.BookmarkService {
	return &bookmarkService{repo: repo}
}

func (s *bookmarkService) CreateBookmark(
	userID uint, rawURL, title, tags string,
) (*domain.Bookmark, error) {
	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return nil, domain.ErrBadRequest("invalid URL format")
	}

	bookmark := &domain.Bookmark{
		UserID: userID,
		URL:    rawURL,
		Title:  title,
		Tags:   tags,
	}

	if err := s.repo.Create(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

func (s *bookmarkService) GetUserBookmarks(
	userID uint,
) ([]domain.Bookmark, error) {
	return s.repo.FindByUserID(userID)
}

func (s *bookmarkService) UpdateBookmark(
	id, userID uint, rawURL, title, tags string,
) (*domain.Bookmark, error) {
	bookmark, err := s.repo.FindByID(id)
	if err != nil {
		return nil, domain.ErrNotFound("bookmark not found")
	}

	if bookmark.UserID != userID {
		return nil, domain.ErrForbidden("you don't own this bookmark")
	}

	if rawURL != "" {
		if _, err := url.ParseRequestURI(rawURL); err != nil {
			return nil, domain.ErrBadRequest("invalid URL format")
		}
		bookmark.URL = rawURL
	}
	if title != "" {
		bookmark.Title = title
	}
	if tags != "" {
		bookmark.Tags = tags
	}

	if err := s.repo.Update(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

func (s *bookmarkService) DeleteBookmark(id, userID uint) error {
	return s.repo.Delete(id, userID)
}