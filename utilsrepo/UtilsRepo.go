// UtilsRepo
package utilsrepo

type (
	UtilsRepo interface {
		CheckUserPassword(userId string, password string) error
		CheckIfUserExists(userId string) error
		CheckIfMatchExists(firstUserId string, secondUserId string) error
		GetUserUsername(userId string) (string, error)
		CheckLoginCredentials(email string, password string) (string, error)
	}
	
	repoFactory func() UtilsRepo
)

var (
	New repoFactory
)
