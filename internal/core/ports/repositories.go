package ports


import (
	"ocb.amot.io/internal/core/domain"
)

type TokenRepoInterface interface {
	FindToken(string)  (*domain.RefreshToken, error)
	FindTokenByUserId(uint)  (*domain.RefreshToken, error)
	UpdateToken(*domain.RefreshToken)  error
	RegisterToken(*domain.RefreshToken)  error
}
