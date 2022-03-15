package userHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"tinder_users/internal/data/Infrastructure/userRepository"
	"tinder_users/internal/encoder"
	"tinder_users/pkg/Domain/response"
	"tinder_users/pkg/Domain/user"

	"github.com/gofrs/uuid"
)

type UserHandler struct {
	Repository userRepository.Repository
}

type Handler interface {
	RegisUser(regisUser user.RegisterUser) response.Status
	Login(loginUser user.RegisterUser) (interface{}, response.Status)
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

func (uh *UserHandler) Login(loginUser user.RegisterUser) (interface{}, response.Status) {
	userFound, status := uh.Repository.FindUserByEmail(loginUser.Email)
	if status != response.UserFound {
		return nil, status
	}

	if encoder.ComparePasswords(userFound.HashedPassword, []byte(loginUser.Password)) {
		userFound.HashedPassword = ""

		resp, err := http.Get("http://localhost:3000/api/tinder/pets/owner/" + userFound.Token)

		if err != nil {
			fmt.Println(err)
			return nil, response.InternalServerError
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {

			var parseResp response.Response
			err = json.NewDecoder(resp.Body).Decode(&parseResp)
			if err != nil {
				return nil, response.InternalServerError
			}

			jsonString, _ := json.Marshal(parseResp.Data)
			err = json.Unmarshal(jsonString, &userFound.Pets)
			if err != nil {
				return nil, response.InternalServerError
			}
		}

		return userFound, response.SuccesfulLogin
	}

	return nil, response.IncorrectPassword
}
