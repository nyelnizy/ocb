package domain

type User struct {
	Id        uint   `json:"id" pact:"example=1"`
	FirstName string `json:"first_name" pact:"example=Daniel"`
	LastName  string `json:"last_name" pact:"example=Addae"`
	Role      string `json:"role" pact:"example=consumer"`
	City      string `json:"city" pact:"example=Accra"`
	Town      string `json:"town" pact:"example=Ablekuma"`
	Email     string `json:"email" pact:"example=yhiamdan@gmail.com"`
}
