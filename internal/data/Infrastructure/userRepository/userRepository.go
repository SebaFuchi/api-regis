package userRepository

import (
	"mascotas_users/pkg/Domain/response"
	"mascotas_users/pkg/Domain/user"
	"mascotas_users/pkg/Use_cases/Helpers/dbHelper"
)

type Repository interface {
	RegisUser(u *user.User) response.Status
	FindUserByEmail(email string) (user.User, response.Status)
}

type UserRepository struct {
}

func (ur *UserRepository) RegisUser(u *user.User) response.Status {
	sqlCon, err := dbHelper.GetDB()
	if err != nil {
		return response.InternalServerError
	}
	defer sqlCon.Close()

	insForm, err := sqlCon.Prepare("INSERT INTO users (token, hashed_password, name, last_name, email)VALUES (?,?,?,?,?)")
	if err != nil {
		return response.DBQueryError
	}
	defer insForm.Close()

	result, err := insForm.Exec(
		u.Token,
		u.HashedPassword,
		u.Name,
		u.LastName,
		u.Email,
	)
	if err != nil {
		return response.DBExecutionError

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

func (ur *UserRepository) FindUserByEmail(email string) (user.User, response.Status) {
	var u user.User
	sqlCon, err := dbHelper.GetDB()
	if err != nil {
		return u, response.InternalServerError
	}
	defer sqlCon.Close()

	selForm, err := sqlCon.Prepare("SELECT id, token, hashed_password, name, last_name, email FROM users WHERE email = ?")
	if err != nil {
		return u, response.DBQueryError
	}
	defer selForm.Close()

	result, err := selForm.Query(
		email,
	)
	if err != nil {
		return u, response.DBExecutionError
	}

	if result.Next() {
		err := result.Scan(
			&u.Id,
			&u.Token,
			&u.HashedPassword,
			&u.Name,
			&u.LastName,
			&u.Email,
		)
		if err != nil {
			return u, response.DBScanError
		}
		return u, response.UserFound
	}

	return u, response.UserNotFound

}
