package userRepository

import (
	"mascotas_users/pkg/Domain/response"
	"mascotas_users/pkg/Domain/user"
	"mascotas_users/pkg/Use_cases/Helpers/dbHelper"
	"regexp"

	"github.com/go-sql-driver/mysql"
)

type Repository interface {
	RegisUser(u *user.User) response.Status
	AddPermit(userId int, from int) response.Status
	FindUser(userTag string) (user.User, response.Status)
	GetUserPermits(userId int) ([]int, response.Status)
}

type UserRepository struct {
}

func (ur *UserRepository) RegisUser(u *user.User) response.Status {
	sqlCon, err := dbHelper.GetDB()
	if err != nil {
		return response.InternalServerError
	}
	defer sqlCon.Close()

	insForm, err := sqlCon.Prepare("INSERT INTO users (email, hashedPassword, accessToken, name, nameTag, regisDate)VALUES (?,?,?,?,?,?)")
	if err != nil {
		return response.DBQueryError
	}
	defer insForm.Close()

	result, err := insForm.Exec(
		u.Email,
		u.HashedPassword,
		u.AccessToken,
		u.Name,
		u.NameTag,
		u.RegisDate,
	)
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if !ok {
			return response.DBExecutionError
		} else {
			if me.Number == 1062 {
				matchEmail, _ := regexp.MatchString("email_UNIQUE", me.Message)
				if matchEmail {
					return response.EmailAlreadyExists
				} else {
					return response.NickNameAlreadyExists
				}
			}
			return response.DBExecutionError
		}
	} else {
		rows, err := result.RowsAffected()
		if err != nil {
			return response.DBRowsError
		}
		if rows == 0 {
			return response.CreationFailure
		}
		id, err := result.LastInsertId()
		if err != nil {
			return response.DBLastRowIdError
		}
		u.Id = int(id)
	}
	return response.SuccesfulCreation
}

func (ur *UserRepository) AddPermit(userId int, from int) response.Status {
	sqlCon, err := dbHelper.GetDB()
	if err != nil {
		return response.InternalServerError
	}
	defer sqlCon.Close()

	insForm, err := sqlCon.Prepare("INSERT INTO user_accessiblelevels (userId,permitId)VALUES (?,?)")
	if err != nil {
		message := "path: /api/hermes-user/users. AddPermit, DBQueryError"
		logger.Log(message, "500")

		return response.DBQueryError
	}
	defer insForm.Close()

	result, err := insForm.Exec(
		userId,
		from,
	)
	if err != nil {
		message := "path: /api/hermes-user/users. AddPermit, DBExecutionError"
		logger.Log(message, "500")

		return response.DBExecutionError
	} else {
		rows, err := result.RowsAffected()
		if err != nil {
			message := "path: /api/hermes-user/users. AddPermit, DBRowsError"
			logger.Log(message, "500")

			return response.DBRowsError
		}
		if rows == 0 {
			message := "path: /api/hermes-user/users. AddPermit, CreationFailure"
			logger.Log(message, "500")

			return response.CreationFailure
		}
	}
	return response.SuccesfulCreation
}

func (ur *UserRepository) FindUser(userTag string) (user.User, response.Status) {
	var u user.User
	sqlCon, err := dbHelper.GetDB()
	if err != nil {
		return u, response.InternalServerError
	}
	defer sqlCon.Close()

	selForm, err := sqlCon.Prepare("SELECT id, email, hashedPassword, accessToken, name, nameTag, regisDate FROM users WHERE email = ?")
	if err != nil {
		return u, response.DBQueryError
	}
	defer selForm.Close()

	err = selForm.QueryRow(userTag, userTag).Scan(
		&u.Id,
		&u.Email,
		&u.HashedPassword,
		&u.AccessToken,
		&u.Name,
		&u.NameTag,
		&u.RegisDate,
	)

	if err != nil {
		return u, response.UserNotFound
	}

	return u, response.UserFound
}

func (ur *UserRepository) GetUserPermits(userId int) ([]int, response.Status) {
	sqlCon, err := dbHelper.GetDB()
	if err != nil {
		return nil, response.InternalServerError
	}
	defer sqlCon.Close()

	selForm, err := sqlCon.Prepare("SELECT permitId FROM user_accessiblelevels WHERE userId = ?;")
	if err != nil {
		message := "path: /hermes-user/api/users. GetUserPermits, preparar sql error"
		logger.Log(message, "500")
		return nil, response.InternalServerError
	}
	var permits []int
	result, err := selForm.Query(userId)
	if err != nil {
		message := "path: /hermes-user/api/users. GetUserPermits, DB ERROR"
		logger.Log(message, "500")
		return nil, response.InternalServerError
	}
	for result.Next() {
		var permitId int
		err := result.Scan(&permitId)
		if err != nil {
			message := "path: /hermes-user/api/users. GetUserPermits, DB ERROR"
			logger.Log(message, "500")
			return nil, response.InternalServerError
		}
		permits = append(permits, permitId)
	}
	return permits, response.UserFound

}
