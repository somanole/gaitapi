// UserRepo
package userrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	UserRepo interface {
		GetUser(userId string) (types.User, error)
		CreateUser(user types.User) (types.CreateUserResponse, error)
		UpdateUser(user types.User) (types.CreateUserResponse, error)
		GetUserByEmail(email string) (types.UserByEmail, error)
		GetUserExtraInfo(userId string) (types.UserExtraInfo, error)
		UpdateUserExtraInfo(ue types.UserExtraInfo) error
		ReportUser(ur types.UserReport) error
	}
	
	repoFactory func() UserRepo
)

var (
	New repoFactory
)
