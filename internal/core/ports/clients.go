package ports

import "ocb.amot.io/internal/core/domain"

type UserClientInterface interface {
	GetUser(string,string) (*domain.User,error)
	GetUserById(uint) (*domain.User,error)
}
