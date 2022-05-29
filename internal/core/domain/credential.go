package domain

type Credential struct {
	Email string `valid:"required,email"`
	Password string `valid:"required"`
}
