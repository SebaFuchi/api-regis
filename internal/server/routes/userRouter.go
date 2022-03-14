package routes

import (
	"encoding/json"
	"hermes_users/pkg/Domain/response"
	"hermes_users/pkg/Domain/user"
	"hermes_users/pkg/Use_cases/Handlers/userHandler"
	"hermes_users/pkg/Use_cases/Helpers/logger"
	permitshelper "hermes_users/pkg/Use_cases/Helpers/permitsHelper"
	"hermes_users/pkg/Use_cases/Helpers/responseHelper"
	"net/http"

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

	registeredUser, status := ur.Handler.RegisUser(u, from)
	resp, err := responseHelper.ResponseBuilder(status.Index(), status.String(), registeredUser)
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
	case response.InvalidEmailFormat:
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(resp))
		return
	case response.EmailAlreadyExists, response.NickNameAlreadyExists:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(resp))
		return
	case response.InternalServerError, response.RequestError, response.RequestDoError, response.ReadRequestError, response.DBQueryError, response.DBExecutionError, response.CreationFailure, response.DBLastRowIdError:
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
	from := chi.URLParam(r, "from")
	var lr user.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&lr)
	if err != nil {
		message := "path: /hermes-user/api/users. Se intenta loguear en hermes, Bad request"
		logger.Log(message, "400")
		resp, err := responseHelper.ResponseBuilder(response.BadRequest.Index(), response.BadRequest.String(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal Server Error"))
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp))
		return
	}

	permit := permitshelper.GetPermit(from)
	if permit == permitshelper.Unknown.Index() {
		resp, err := responseHelper.ResponseBuilder(response.OriginNotAllowed.Index(), response.OriginNotAllowed.String(), nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: Internal server error"))
			return
		}
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(resp))
		return
	}
	defer r.Body.Close()

	userParsed, status := ur.Handler.Login(lr, from)
	var resp []byte
	resp, err = responseHelper.ResponseBuilder(status.Index(), status.String(), userParsed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: Internal Server Error"))
		return
	}
	switch status {
	// There is no user with that data to login
	case response.UserDontExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(resp))
		return
	// the password does not match the one found in the token
	case response.IncorrectPassword:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(resp))
		return
	case response.InvalidPermissions:
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(resp))
		return
	case response.UserFound:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
		return
	case response.RequestTimeOut:
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte(resp))
	default:
		message := "path: /hermes-user/api/users. Se intenta loguear en hermes, ERROR INESPERADO."
		logger.Log(message, "500")
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

	//r.Post("/login", ur.Login)
	r.Post("/register/{from}", ur.RegisUser)
	r.Post("/login/{from}", ur.Login)

	return r
}
