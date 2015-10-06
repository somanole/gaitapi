// MatchRepo
package matchrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	MatchRepo interface {
		CreateMatch(userId string, matchedUserId string) error
		GetUserMatch(userId string) (types.Match, error)
	}
	
	repoFactory func() MatchRepo
)

var (
	New repoFactory
)
