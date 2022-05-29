package controllers

import (
	"net/http"
	"ocb.amot.io/internal/core/domain"
)

func generateError(w http.ResponseWriter,err error,code int){
	httpErr := domain.HttpErr{
		Err:    err,
		Status: code,
	}
	http.Error(w, httpErr.String() ,code)
}
