// UtilsRepo
package utilsrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	UtilsRepo interface {
		CheckUserPassword(userId string, password string) error
		CheckIfUserExists(userId string) error
		CheckIfMatchExists(firstUserId string, secondUserId string) error
		GetUserUsername(userId string) (string, error)
		CheckLoginCredentials(email string, password string) (string, error)
		RegisterInterest(i types.Interest) error
	}
	
	repoFactory func() UtilsRepo
)

var (
	New repoFactory
)
