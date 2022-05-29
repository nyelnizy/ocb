package contract

import (
	"fmt"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/suite"
	"net/http"
	"ocb.amot.io/internal/adapters/clients"
	"ocb.amot.io/internal/core/domain"
	"testing"
)


type UserClientSuite struct {
  suite.Suite
}

var (
	pact *dsl.Pact
)

func TestUserClientSuite(t *testing.T)  {
	suite.Run(t,new(UserClientSuite))
}
func (s *UserClientSuite) SetupSuite()  {
	pact = &dsl.Pact{
		Consumer: "AuthService",
		Provider: "UserService",
		PactDir: "./pacts",
	}
}

func (s *UserClientSuite) TestGetUser() {
	// add interactions, defines what we intend
	// to request and what we expect as response
	pact.AddInteraction().
		Given("There is a user with email yhiamdan@gmail.com and password 1234").
		UponReceiving("A GET request for a user with email yhiamdan@gmail.com and password 1234").
		WithRequest(dsl.Request{
		Method:  http.MethodGet,
		Path:    dsl.String("/api/users"),
		Query:   dsl.MapMatcher{
			"email": dsl.String("yhiamdan@gmail.com"),
			"password": dsl.String("1234"),
		},
		Headers: dsl.MapMatcher{
			"Accept":dsl.String("application/json"),
		},
	}).WillRespondWith(dsl.Response{
		Status:  http.StatusOK,
		Headers: dsl.MapMatcher{
			"Content-Type":dsl.String("application/json"),
		},
		Body:    dsl.Match(domain.User{}),
	})

	// sends a dummy fake request to a mock user services and verifies response
	test := func() error{
		c := clients.NewUserClient(fmt.Sprintf("http://localhost:%d",pact.Server.Port))
		u, err := c.GetUser("yhiamdan@gmail.com","1234")
		s.NoError(err)
		s.Equal(uint(1),u.Id)
		s.Equal("Daniel",u.FirstName)
		s.Equal("Addae",u.LastName)
		s.Equal("consumer",u.Role)
		s.Equal("Accra",u.City)
		s.Equal("Ablekuma",u.Town)
		s.Equal("yhiamdan@gmail.com",u.Email)
		return nil
	}
	s.NoError(pact.Verify(test))
}

func (s *UserClientSuite) TestGetUserById() {
	// add interactions, defines what we intend
	// to request and what we expect as response
	pact.AddInteraction().
		Given("There is a user with id 1").
		UponReceiving("A GET request for a user with id 1").
		WithRequest(dsl.Request{
			Method:  http.MethodGet,
			Path:    dsl.String("/api/users"),
			Query: dsl.MapMatcher{
				"id": dsl.String("1"),
			},
			Headers: dsl.MapMatcher{
				"Accept":dsl.String("application/json"),
			},
		}).WillRespondWith(dsl.Response{
		Status:  http.StatusOK,
		Headers: dsl.MapMatcher{
			"Content-Type":dsl.String("application/json"),
		},
		Body:    dsl.Match(domain.User{}),
	})

	// sends a dummy fake request to a mock user services and verifies response
	test := func() error{
		c := clients.NewUserClient(fmt.Sprintf("http://localhost:%d",pact.Server.Port))
		u, err := c.GetUserById(uint(1))
		s.NoError(err)
		s.Equal(uint(1),u.Id)
		s.Equal("Daniel",u.FirstName)
		s.Equal("Addae",u.LastName)
		s.Equal("consumer",u.Role)
		s.Equal("Accra",u.City)
		s.Equal("Ablekuma",u.Town)
		s.Equal("yhiamdan@gmail.com",u.Email)
		return nil
	}
	s.NoError(pact.Verify(test))
}

func (s *UserClientSuite) TearDownSuite()  {
	pact.Teardown()
}