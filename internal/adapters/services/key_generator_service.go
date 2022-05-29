package services

import (
	uuid "github.com/nu7hatch/gouuid"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type KeyGeneratorService struct {
	path string
}

func NewKeyGenerator(path string) *KeyGeneratorService {
	return &KeyGeneratorService{path}
}

func (k *KeyGeneratorService) GenerateRsaKey() ([]byte,error) {
	relPath, err := filepath.Abs(k.path)
	if err!=nil{
		log.Fatal(err)
		return nil, err
	}
	keyContents,err := ioutil.ReadFile(relPath)
	if err!=nil{
		log.Fatal(err)
		return nil, err
	}
	return keyContents,nil
}
func (k *KeyGeneratorService) GenerateRefreshToken() (string, error)  {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return uid.String(), nil
}

func (k *KeyGeneratorService) GenerateExpirationTime(minutes int) int64  {
	return time.Now().Add(time.Minute*time.Duration(minutes)).Unix()
}