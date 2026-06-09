package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/rwndy/bookmark-api/internal/domain"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthService struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	jwtSecret        []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
}

func NewAuthService(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	jwtSecret string,
	accessTokenTTL, refreshTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        []byte(jwtSecret),
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
	}
}

func (s *AuthService) Register(email, password string) (*domain.User, error) {
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, domain.ErrBadRequest("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashed),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, domain.ErrUnauthorized("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrUnauthorized("invalid email or password")
	}

	return s.issueTokenPair(user.ID, user.Email)
}

func (s *AuthService) Refresh(refreshToken string) (*TokenPair, error) {
	if refreshToken == "" {
		return nil, domain.ErrUnauthorized("invalid refresh token")
	}

	hash := hashToken(refreshToken)
	stored, err := s.refreshTokenRepo.FindByTokenHash(hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUnauthorized("invalid refresh token")
		}
		return nil, err
	}

	if !stored.IsActive(time.Now()) {
		return nil, domain.ErrUnauthorized("refresh token expired or revoked")
	}

	user, err := s.userRepo.FindByID(stored.UserID)
	if err != nil {
		return nil, domain.ErrUnauthorized("invalid refresh token")
	}

	// Rotate: revoke the old refresh token before issuing a new pair.
	if err := s.refreshTokenRepo.Revoke(stored.ID); err != nil {
		return nil, err
	}

	return s.issueTokenPair(user.ID, user.Email)
}

func (s *AuthService) Logout(refreshToken string) error {
	if refreshToken == "" {
		return domain.ErrBadRequest("refresh token is required")
	}

	stored, err := s.refreshTokenRepo.FindByTokenHash(hashToken(refreshToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return s.refreshTokenRepo.Revoke(stored.ID)
}

func (s *AuthService) issueTokenPair(userID uint, email string) (*TokenPair, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     now.Add(s.accessTokenTTL).Unix(),
		"iat":     now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshRaw, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	stored := &domain.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(refreshRaw),
		ExpiresAt: now.Add(s.refreshTokenTTL),
	}
	if err := s.refreshTokenRepo.Create(stored); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshRaw,
		ExpiresIn:    int64(s.accessTokenTTL.Seconds()),
	}, nil
}

func generateRandomToken(numBytes int) (string, error) {
	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
