// MatchRepo
package matchrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	MatchRepo interface {
		GetUserMatch(userId string) (types.Match, error)
	}
	
	repoFactory func() MatchRepo
)

var (
	New repoFactory
)
