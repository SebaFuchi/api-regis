package userHandler

import (
	"mascotas_users/internal/data/Infrastructure/userRepository"
	"mascotas_users/internal/encoder"
	"mascotas_users/pkg/Domain/response"
	"mascotas_users/pkg/Domain/user"
	"net/mail"

	"github.com/gofrs/uuid"
)

type UserHandler struct {
	Repository userRepository.Repository
}

type Handler interface {
	RegisUser(regisUser user.RegisterUser) response.Status
	Login(loginUser user.LoginRequest) (interface{}, response.Status)
}

func (ur *UserHandler) RegisUser(regisUser user.RegisterUser) response.Status {

	_, status := ur.Repository.FindUserByEmail(regisUser.Email)
	if status == response.UserFound {
		return response.EmailAlreadyExists
	}

	token, err := uuid.NewV4()
	if err != nil {
		return response.InternalServerError
	}
	_, err = mail.ParseAddress(regisUser.Email)
	if err != nil {
		return response.InvalidEmailFormat
	}

	u := user.User{
		Token:          token.String(),
		Email:          regisUser.Email,
		Name:           regisUser.Name,
		LastName:       regisUser.LastName,
		HashedPassword: string(encoder.HashAndSalt([]byte(regisUser.Password))),
	}

	status = ur.Repository.RegisUser(&u)
	if status != response.SuccesfulCreation {
		return status
	}

	return response.SuccesfulCreation
}

func (uh *UserHandler) Login(loginUser user.LoginRequest) (interface{}, response.Status) {
	userFound, status := uh.Repository.FindUserByEmail(loginUser.Email)
	if status != response.UserFound {
		return nil, status
	}

	if encoder.ComparePasswords(userFound.HashedPassword, []byte(loginUser.Password)) {
		userFound.HashedPassword = ""
		return userFound, response.SuccesfulLogin
	}
	return nil, response.IncorrectPassword
}
