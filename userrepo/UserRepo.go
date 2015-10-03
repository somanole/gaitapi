// UserRepo
package userrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	UserRepo interface {
		GetUser(userId string) (types.User, error)
		CreateUser(user types.UserRequest) (types.CreateUserResponse, error)
		UpdateUser(userId string, user types.UserUpdateRequest) (types.CreateUserResponse, error)
		GetUserByEmail(email string) (types.UserByEmail, error)
		GetUserExtraInfo(userId string) (types.UserExtraInfo, error)
	}
	
	repoFactory func() UserRepo
)

var (
	New repoFactory
)
