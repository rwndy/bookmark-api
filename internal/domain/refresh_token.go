package domain

import "time"

type RefreshToken struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"not null;index"`
	TokenHash string     `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (t *RefreshToken) IsActive(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	FindByTokenHash(hash string) (*RefreshToken, error)
	Revoke(id uint) error
	RevokeAllByUserID(userID uint) error
}
