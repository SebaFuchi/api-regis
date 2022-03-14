package routes

import (
	"mascotas_users/internal/data/Infrastructure/userRepository"
	"mascotas_users/pkg/Use_cases/Handlers/userHandler"

	"net/http"

	"github.com/go-chi/chi"
)

// Instanciamos los handlers de los endpoints
func New() http.Handler {
	r := chi.NewRouter()

	ur := &UserRouter{
		Handler: &userHandler.UserHandler{
			Repository: &userRepository.UserRepository{},
		},
	}

	r.Mount("/users", ur.Routes())

	//Retornamos la api ya construida
	return r

}
