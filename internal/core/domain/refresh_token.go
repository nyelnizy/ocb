package domain

import "gorm.io/gorm"

type RefreshToken struct {
	gorm.Model
	UserId               uint   `json:"user_id"`
	Token                string `json:"token"`
	Invalidated          bool   `json:"invalidated"`
	ExpiresAt            int64  `json:"expires_at"`
	AccessTokenExpiresAt int64  `json:"access_token_expires_at"`
}
