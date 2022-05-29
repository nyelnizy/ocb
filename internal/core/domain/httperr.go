package domain

import "os"

type HttpErr struct {
	Err error
	Status int
}

func (h *HttpErr) String() string{
	env := os.Getenv("APP_ENV")
	if h.Status==500 && env=="prod"{
		return "Internal Server Error"
	}
	return h.Err.Error()
}