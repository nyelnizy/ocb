package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ocb.amot.io/internal/core/domain"
)

type UserClient struct {
 host string
}

func NewUserClient(host string)  *UserClient {
	return &UserClient{host: host}
}

func (u *UserClient) GetUser(email string,password string) (*domain.User,error) {
	return getUser(fmt.Sprintf("%s/api/users?email=%s&password=%s",u.host,email,password))
}

func (u *UserClient) GetUserById(userId uint) (*domain.User, error) {
	return getUser(fmt.Sprintf("%s/api/users?id=%d",u.host,userId))
}

func getUser(path string) (*domain.User, error) {
	client := http.Client{}
	req, _ := http.NewRequest(
		http.MethodGet,
		path,
		nil,
	)
	req.Header.Add("Accept","application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	us := &domain.User{}
	json.NewDecoder(resp.Body).Decode(us)
	return us, nil
}