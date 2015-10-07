// ActivityRepo
package activityrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	ActivityRepo interface {
		CreateUserActivity(a types.Activity) error
		GetUserActivity(userId string) (types.Activity, error)
	}
	
	repoFactory func() ActivityRepo
)

var (
	New repoFactory
)
