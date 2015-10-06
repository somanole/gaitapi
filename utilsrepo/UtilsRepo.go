// UtilsRepo
package utilsrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	UtilsRepo interface {
		CheckUserPassword(userId string, rpassword string) error
		CheckIfUserExists(userId string) error
		GetUserUsername(userId string) string, error
	}
	
	repoFactory func() UtilsRepo
)

var (
	New repoFactory
)
