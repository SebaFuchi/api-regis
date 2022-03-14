package userHandler

import (
	"fmt"
	"mascotas_users/internal/data/Infrastructure/userRepository"
	"mascotas_users/internal/encoder"
	"mascotas_users/pkg/Domain/response"
	"mascotas_users/pkg/Domain/user"
	"mascotas_users/pkg/Use_cases/Helpers/logger"
	"net/mail"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type UserHandler struct {
	Repository userRepository.Repository
}

type Handler interface {
	RegisUser(regisUser user.RegisterUser) (interface{}, response.Status)
	Login(loginUser user.LoginRequest, from string) (interface{}, response.Status)
}

func (ur *UserHandler) RegisUser(regisUser user.RegisterUser) (interface{}, response.Status) {

	token, err := uuid.NewV4()
	if err != nil {
		return nil, response.InternalServerError
	}
	_, err = mail.ParseAddress(regisUser.Email)
	if err != nil {
		return nil, response.InvalidEmailFormat
	}

	u := user.User{
		Token:          token.String(),
		Email:          regisUser.Email,
		RegisDate:      fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
		Name:           fmt.Sprintf("%s %s", regisUser.Name, regisUser.LastName),
		NameTag:        regisUser.NameTag,
		HashedPassword: string(encoder.HashAndSalt([]byte(regisUser.Password))),
	}

	status := ur.Repository.RegisUser(&u)
	if status != response.SuccesfulCreation {
		return nil, status
	}

	return aegirHelper.RegisAegirUser(u)

}

func (uh *UserHandler) Login(loginUser user.LoginRequest, from string) (interface{}, response.Status) {
	var userFound user.User
	var status response.Status

	// manejar logueo con email o con nickname

	userFound, status = uh.Repository.FindUser(loginUser.User)
	if status != response.UserFound {
		return nil, status
	}

	if encoder.ComparePasswords(userFound.HashedPassword, []byte(loginUser.Password)) {
		sessId, err := uuid.NewV4()
		if err != nil {
			message := "path: /hermes-user/api/users. Se intenta comparar contraseñas para loguear, ERROR"
			logger.Log(message, "500")
			return nil, response.InternalServerError
		}
		userFound.SessionId = sessId.String()
		permitIds, status := uh.Repository.GetUserPermits(userFound.Id)
		if status != response.UserFound {
			return nil, status
		}

		for _, permit := range permitIds {
			if permit == permitshelper.GetPermit(from) {
				destiny := strings.Split(from, "-")
				switch destiny[0] {
				case "aegir":
					return aegirHelper.LoginAegir(userFound.Email)
				}
			}
		}
		return nil, response.InvalidPermissions
	}
	return nil, response.IncorrectPassword
}
