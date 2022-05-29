package unit

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"ocb.amot.io/internal/adapters/services"
	"ocb.amot.io/internal/core/domain"
	"testing"
	"time"
)


var id uint = 1
var expirationTime int64= 1652093130109
var refreshToken = "xxxxxxxxx"
var role = "consumer"

func TestIssueToken(t *testing.T) {
	// create instance of mock repo
	tokenRepo := new(MockedTokenRepo)

	// setup expectations
	tokenRepo.On("RegisterToken", &domain.RefreshToken{
		UserId:               id,
		Token:                refreshToken,
		Invalidated:          false,
		ExpiresAt:            expirationTime,
		AccessTokenExpiresAt: expirationTime,
	}).Return( nil)

	// create instance of mock user entry points
	userService := new(MockedUserService)
	// setup expectations
	u := domain.User{
		Id:        id,
		FirstName: "Daniel",
		LastName:  "Addae",
		Role:      role,
		City:      "Accra",
		Town:      "Ablekuma",
		Email:     "yhiamdan@gmail.com",
	}
	userService.On("GetUser","yhiamdan@gmail.com","password").Return(&u,nil)

	// create instance of mocked key generator
	keyGenerator := new(MockedKeyGenerator)
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// setup expectations
	keyGenerator.On("GenerateRsaKey").Return(keyPEM,nil)
	keyGenerator.On("GenerateRefreshToken").Return(refreshToken,nil)
	keyGenerator.On("GenerateExpirationTime",30).Return(expirationTime)
	keyGenerator.On("GenerateExpirationTime",60*60).Return(expirationTime)

	// test
	ts := services.NewTokenService(tokenRepo,userService,keyGenerator)
    token, err := ts.IssueToken(&domain.Credential{
		Email:    "yhiamdan@gmail.com",
		Password: "password",
	})

    // assert
    tokenRepo.AssertExpectations(t)
	userService.AssertExpectations(t)
	keyGenerator.AssertExpectations(t)

    assert.Nil(t,err)
    assert.NotNil(t,token)
    assert.NotNil(t,token.Access)
    assert.NotNil(t,token.Refresh)
}

func TestVerifyToken(t *testing.T)  {
	// create instance of mocked key generator
	keyGenerator := new(MockedKeyGenerator)

	// setup expectations
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	keyGenerator.On("GenerateRsaKey").Return(keyPEM,nil)

	// prepare
	tc := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), services.CustomClaims{
		UserId: id,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
			Issuer:    services.Issuer,
		},
	})

	token, err := tc.SignedString(key)
	if err !=nil{
		panic(err)
	}

	// test
	ts := services.NewTokenService(nil,nil,keyGenerator)
	parsedToken, err := ts.VerifyToken(token)

	// assert
	assert.Nil(t,err)
	assert.NotNil(t,parsedToken)

	//verify claims
	claims := parsedToken.Claims.(*services.CustomClaims)
	assert.Equal(t, id,claims.UserId)
	assert.Equal(t, services.Issuer,claims.Issuer)
	assert.Equal(t,"consumer",claims.Role)
}

func TestInvalidateToken(t *testing.T)  {
	rt := &domain.RefreshToken{UserId: id}
	tokenRepo := new(MockedTokenRepo)

	// setup expectations
	tokenRepo.On("FindTokenByUserId", id).Return(rt, nil)
	tokenRepo.On("UpdateToken", rt).Return(nil)

	// test
	ts := services.NewTokenService(tokenRepo,nil,nil)
	err := ts.InvalidateToken(id)

	// assert
	tokenRepo.AssertExpectations(t)
	assert.Nil(t,err)
}

func TestRefreshToken(t *testing.T)  {
	// must be expired
	accessTime := time.Now().Add(-20*time.Minute)

	// must be valid
	refreshTime := time.Now().Add(time.Minute*20)

	rt := &domain.RefreshToken{UserId: id,Token: refreshToken,AccessTokenExpiresAt: accessTime.Unix(),ExpiresAt:refreshTime.Unix() }
	tokenRepo := new(MockedTokenRepo)

	tokenRepo.On("RegisterToken", rt).Return( nil)
	tokenRepo.On("FindToken", refreshToken).Return(rt, nil)

	// create instance of mock user entry points
	userService := new(MockedUserService)

	// setup expectations
	u := domain.User{
		Id:        id,
		FirstName: "Daniel",
		LastName:  "Addae",
		Role:      role,
		City:      "Accra",
		Town:      "Ablekuma",
		Email:     "yhiamdan@gmail.com",
	}
	userService.On("GetUserById",rt.UserId).Return(&u,nil)

	// create instance of mocked key generator
	keyGenerator := new(MockedKeyGenerator)

	// setup expectations
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	keyGenerator.On("GenerateRsaKey").Return(keyPEM,nil)
	keyGenerator.On("GenerateRefreshToken").Return(refreshToken,nil)
	keyGenerator.On("GenerateExpirationTime",30).Return(rt.AccessTokenExpiresAt)
	keyGenerator.On("GenerateExpirationTime",60*60).Return(rt.ExpiresAt)

	ts := services.NewTokenService(tokenRepo,userService,keyGenerator)

	// test
	token, err:= ts.RefreshToken(refreshToken)

	// assert
	tokenRepo.AssertExpectations(t)
	userService.AssertExpectations(t)
	keyGenerator.AssertExpectations(t)

	assert.Nil(t,err)
	assert.NotNil(t,token)
	assert.NotNil(t,token.Access)
	assert.NotNil(t,token.Refresh)
}