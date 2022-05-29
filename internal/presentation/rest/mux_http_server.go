package presentation

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"ocb.amot.io/internal/adapters/services"
	"ocb.amot.io/internal/core/ports"
	"ocb.amot.io/internal/presentation/rest/controllers"
	"strings"
)

var tokenService ports.TokenServiceInterface

type MuxHttpServer struct {
	lc *controllers.LoginController
}

func NewMuxHttpServer(ts ports.TokenServiceInterface) *MuxHttpServer {
	tokenService = ts
	lc := controllers.NewLoginController(ts)
	return &MuxHttpServer{lc}
}

func (h *MuxHttpServer) Start(host string, port string) error {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s:%s/swagger/doc.json", host, port)),
	)).Methods(http.MethodGet)

	r.PathPrefix("/api/v1")
	r.HandleFunc("/login", h.lc.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", requiresAuth(h.lc.Logout)).Methods(http.MethodPost)
	r.HandleFunc("/refresh/{token}", h.lc.RefreshToken).Methods(http.MethodPost)
	log.Printf("Started a new rest entrypoints on %s...", host)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}
	return srv.ListenAndServe()
}

func requiresAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}
		s := strings.Split(token, " ")
		if len(s) != 2 {
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}
		tokenStr := s[1]
		tok, err := tokenService.VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthenticated", http.StatusUnauthorized)
			return
		}
		claims := tok.Claims.(services.CustomClaims)
        ctx := context.WithValue(r.Context(),"userId",claims.UserId)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
