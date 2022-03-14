package routes

import (
	"encoding/json"
	"net/http"
	"tinder_users/pkg/Domain/response"
	"tinder_users/pkg/Domain/user"
	"tinder_users/pkg/Use_cases/Handlers/userHandler"
	"tinder_users/pkg/Use_cases/Helpers/responseHelper"

	"github.com/go-chi/chi"
)

type UserRouter struct {
	Handler userHandler.Handler
}

func (ur *UserRouter) RegisUser(w http.ResponseWriter, r *http.Request) {
	var u user.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		status := response.BadRequest
		response, err := responseHelper.ResponseBuilder(status.Index(), status.String(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal server error"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(response))
		return
	}

	status := ur.Handler.RegisUser(u)
	resp, err := responseHelper.ResponseBuilder(status.Index(), status.String(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal server error"))
		return
	}
	switch status {
	case response.SuccesfulCreation:
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(resp))
		return
	case response.BadRequest:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp))
		return
	case response.EmailAlreadyExists:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(resp))
		return
	case response.InternalServerError, response.DBQueryError, response.DBExecutionError, response.CreationFailure, response.DBLastRowIdError:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(resp))
		return
	default:
		status = response.Unknown
		response, err := responseHelper.ResponseBuilder(status.Index(), status.String(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal server error"))
			return
		}
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(response))
		return
	}
}

func (ur *UserRouter) Login(w http.ResponseWriter, r *http.Request) {
	var lr user.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		resp, err := responseHelper.ResponseBuilder(response.BadRequest.Index(), response.BadRequest.String(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp))
		return
	}
	defer r.Body.Close()

	userParsed, status := ur.Handler.Login(lr)
	var resp []byte
	resp, err = responseHelper.ResponseBuilder(status.Index(), status.String(), userParsed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
	switch status {
	case response.UserNotFound:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(resp))
		return
	case response.IncorrectPassword:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(resp))
		return
	case response.SuccesfulLogin:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
		return
	case response.RequestTimeOut:
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte(resp))
	default:
		resp, err := responseHelper.ResponseBuilder(response.InternalServerError.Index(), response.InternalServerError.String(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Errors"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(resp))
		return
	}
}

func (ur *UserRouter) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", ur.RegisUser)
	r.Post("/login", ur.Login)

	return r
}
