package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"ocb.amot.io/internal/core/domain"
	"ocb.amot.io/internal/core/ports"
	"time"
)

const Issuer = "https://auth.ocb.com"

type CustomClaims struct {
	UserId uint
	Role string
	jwt.StandardClaims
}
type TokenService struct {
	tr ports.TokenRepoInterface
	us ports.UserClientInterface
	k  ports.KeyGeneratorInterface
}
func NewTokenService(tr ports.TokenRepoInterface, us ports.UserClientInterface,k ports.KeyGeneratorInterface) *TokenService {
	return &TokenService{tr, us, k}
}

func (t *TokenService) GenerateToken(u *domain.User) (*domain.Token, error)   {
	// generate refresh token
	refresh, err := t.k.GenerateRefreshToken()
	if err!=nil{
		log.Fatal(err)
		return nil, err
	}

	accessExpires := t.k.GenerateExpirationTime(30)
	refreshExpires := t.k.GenerateExpirationTime(60*60)

	// register refresh token, access token expiration time is required
	// to verify access token has expired before issuing a new token using refresh token
	token := &domain.RefreshToken{
		UserId: u.Id,
		Token:                refresh,
		Invalidated:          false,
		ExpiresAt:            refreshExpires,
		AccessTokenExpiresAt: accessExpires,
	}
	err = t.tr.RegisterToken(token)
	if err !=nil{
		log.Fatal(err)
		return nil, err
	}

	// generate access token
	keyContents, err := t.k.GenerateRsaKey()
	if err!=nil{
		log.Fatal(err)
		return nil, err
	}
	pk, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)

	if err!=nil{
		log.Fatal(err)
		return nil, err
	}
	tc := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), CustomClaims{
		UserId:         u.Id,
		Role: u.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpires,
			IssuedAt:  time.Now().Unix(),
			Issuer:    Issuer,
		},
	})
	access, err := tc.SignedString(pk)
	if err !=nil{
		log.Fatal(err)
		return nil, err
	}

	return &domain.Token{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (t *TokenService) IssueToken(credential *domain.Credential) (*domain.Token, error) {

	// get user details fro user microservice
	u ,err:= t.us.GetUser(credential.Email, credential.Password)
	if err != nil {
		return nil, err
	}
	token, err := t.GenerateToken(u)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t *TokenService) InvalidateToken(userId uint) error {
	token, err := t.tr.FindTokenByUserId(userId)
	if err != nil {
		return err
	}
	token.Invalidated = true

	err =t.tr.UpdateToken(token)
	if err!=nil{
		return err
	}
	return nil
}

func (t *TokenService) RefreshToken(refreshToken string) (*domain.Token, error) {
	// get instance of refresh token
	tkn, err := t.tr.FindToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// to refresh token, access token mush have expired, refresh token must not have expired
	// and should be valid
	if tkn.Invalidated || tkn.AccessTokenExpiresAt>time.Now().Unix() || tkn.ExpiresAt<time.Now().Unix(){
		tkn.Invalidated = true
		_ = t.tr.UpdateToken(tkn)
		return nil, errors.New("token can not be refreshed")
	}

	// get associated user from user microservice
	u ,err:= t.us.GetUserById(tkn.UserId)
	if err != nil {
		return nil, err
	}
	// generate new token
	token, err := t.GenerateToken(u)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t *TokenService) VerifyToken(tokenString string) (*jwt.Token, error) {
	// first obtain rsa private key and then public key for verification
	keyContents, err:= t.k.GenerateRsaKey()
	pk, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// parse supplied token and verify contents
	token, err := jwt.ParseWithClaims(tokenString,&CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return pk.Public(), nil
	})

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("That's not even a token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				fmt.Println("Timing is everything")
			} else {
				fmt.Println("Couldn't handle this token:", err)
			}
		} else {
			fmt.Println("Couldn't handle this token:", err)
		}
		return nil, err
	}
	return token,nil
}