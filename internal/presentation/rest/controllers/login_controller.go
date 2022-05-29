package controllers

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"net/http"
	"ocb.amot.io/internal/core/domain"
	"ocb.amot.io/internal/core/ports"
)

type LoginController struct {
	ts ports.TokenServiceInterface
}

func NewLoginController(ts ports.TokenServiceInterface) *LoginController {
	return &LoginController{ts}
}

// @Summary      Login endpoint
// @Description  Login by submitting email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  domain.Token
// @Failure      400  {string}  string  true  "Bad Request"
// @Failure      422  {string}  string  true  "Unprocessed Entity"
// @Failure      500  {string}  string  true  "Internal Server Error"
// @Router       /api/v1/login [post]
func (c *LoginController) Login(w http.ResponseWriter, r *http.Request) {

	var b []byte
	_, err := r.Body.Read(b)
	if err != nil {
		generateError(w,err,http.StatusBadRequest)
		return
	}
	cred := domain.Credential{}
	if err := json.Unmarshal(b, &cred); err != nil {
		generateError(w,err,http.StatusBadRequest)
		return
	}
	_, err = govalidator.ValidateStruct(cred)
	if err != nil {
		generateError(w,err,http.StatusUnprocessableEntity)
		return
	}
	token, err := c.ts.IssueToken(&cred)
	if err != nil {
		generateError(w,err,http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(token)
	if err != nil {
		generateError(w,err,http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(res)
}

// @Summary      Refresh token endpoint
// @Description  Refresh token by submitting a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  path      string  true  "103ce042-c9ef-451e-a74e-a7e8d36cd3bd"
// @Success      200    {object}  domain.Token
// @Failure      500    {string}  string  true  "Internal Server Error"
// @Router       /api/v1/refresh/{token} [post]
func (c *LoginController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, err := c.ts.RefreshToken(vars["token"])
	if err != nil {
		generateError(w,err,http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(token)
	if err != nil {
		generateError(w,err,http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(res)
}

// @Summary      Logout endpoint
// @Description  An authenticated user can logout using this endpoint
// @Tags         auth
// @Security ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "Logged Out Successfully"
// @Router       /api/v1/logout [post]
func (c *LoginController) Logout(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId").(uint)
	_ = c.ts.InvalidateToken(userId)
	_, _ = w.Write([]byte("Logged Out Successfully"))
	return
}
