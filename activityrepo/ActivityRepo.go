// ActivityRepo
package activityrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	ActivityRepo interface {
		CreateUserActivity(userId string, a types.ActivityRequest) error
		GetUserActivity(userId string) (types.Activity, error)
	}
	
	repoFactory func() ActivityRepo
)

var (
	New repoFactory
)
