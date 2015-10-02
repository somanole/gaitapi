// AccelerationRepo
package accelerationrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	AccelerationRepo interface {
		CreateAcceleration(userId string, acceleration types.AccelerationRequest) error
	}
	
	repoFactory func() AccelerationRepo
)

var (
	New repoFactory
)
