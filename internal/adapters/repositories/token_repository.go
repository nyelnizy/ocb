package repositories

import (
	"gorm.io/gorm"
	"ocb.amot.io/internal/core/domain"
)

type TokenRepository struct {
	db gorm.DB
}

func NewTokenRepository(db *gorm.DB)  *TokenRepository {
	return &TokenRepository{*db}
}
func (t *TokenRepository) FindToken(refreshToken string) (*domain.RefreshToken, error)  {
	token := domain.RefreshToken{Token: refreshToken}
	r := t.db.First(&token)
	return &token,r.Error
}
func (t *TokenRepository) FindTokenByUserId(userId uint) (*domain.RefreshToken, error)  {
	token := domain.RefreshToken{UserId: userId}
	r := t.db.First(&token)
	return &token,r.Error
}
func (t *TokenRepository) UpdateToken(token *domain.RefreshToken) error  {
    r :=t.db.Save(&token)
	return r.Error
}
func (t *TokenRepository) RegisterToken(token *domain.RefreshToken) error  {
	r:=t.db.Create(&token)
	return r.Error
}