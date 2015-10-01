// UserRepo
package userrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	UserRepo interface {
		GetUser(userId string) (types.User, error)
		CreateUser(user types.User) types.CreateUserResponse
		UpdateUser(userId string, u types.User) (types.CreateUserResponse, error)
	}
	
	repoFactory func() UserRepo
)

var (
	New repoFactory
)
