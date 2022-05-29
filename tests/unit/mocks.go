package unit

import (
	"github.com/stretchr/testify/mock"
	"ocb.amot.io/internal/core/domain"
)

type MockedTokenRepo struct{
	mock.Mock
}
func (m *MockedTokenRepo) FindTokenByUserId(userId uint)  (*domain.RefreshToken, error){
	args := m.Called(userId)
	return args[0].(*domain.RefreshToken),args.Error(1)
}
func (m *MockedTokenRepo) FindToken(refreshToken string)  (*domain.RefreshToken, error){
	args := m.Called(refreshToken)
	return args[0].(*domain.RefreshToken),args.Error(1)
}
func (m *MockedTokenRepo) UpdateToken(token *domain.RefreshToken)  error{
	args := m.Called(token)
	return args.Error(0)
}
func (m *MockedTokenRepo) RegisterToken(token *domain.RefreshToken)  error{
	args := m.Called(token)
	return args.Error(0)
}

type MockedUserService struct {
	mock.Mock
}

func (m *MockedUserService) GetUser(email string,password string) (*domain.User,error)  {
	args := m.Called(email,password)
	return args[0].(*domain.User),args.Error(1)
}

func (m *MockedUserService) GetUserById(userId uint) (*domain.User,error)  {
	args := m.Called(userId)
	return args[0].(*domain.User),args.Error(1)
}

type MockedKeyGenerator struct {
	mock.Mock
}

func (m *MockedKeyGenerator) GenerateRsaKey() ([]byte, error) {
	args := m.Called()
	return args[0].([]byte),args.Error(1)
}
func (m *MockedKeyGenerator) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0),args.Error(1)
}
func (m *MockedKeyGenerator) GenerateExpirationTime(minutes int) int64  {
	args := m.Called(minutes)
	return args[0].(int64)
}

