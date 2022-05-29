package ports

import (
	"github.com/golang-jwt/jwt"
	"ocb.amot.io/internal/core/domain"
)

type TokenServiceInterface interface {
	GenerateToken(*domain.User) (*domain.Token, error)
	IssueToken(*domain.Credential) (*domain.Token, error)
	InvalidateToken(uint) error
	RefreshToken(string) (*domain.Token, error)
	VerifyToken(string) (*jwt.Token,error)
}

type KeyGeneratorInterface interface {
	GenerateRsaKey() ([]byte,error)
	GenerateRefreshToken() (string,error)
	GenerateExpirationTime(int) int64
}

