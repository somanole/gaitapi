// MatchRepo
package matchrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	MatchRepo interface {
		CreateMatch(m types.Match) error
		GetUserMatch(userId string) (types.Match, error) 
	}
	
	repoFactory func() MatchRepo
)

var (
	New repoFactory
)
