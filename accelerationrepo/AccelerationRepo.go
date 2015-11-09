// AccelerationRepo
package accelerationrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	AccelerationRepo interface {
		CreateAcceleration(a types.Acceleration) error
		GetAccelerations(userId string) (types.Accelerations, error)
	}
	
	repoFactory func() AccelerationRepo
)

var (
	New repoFactory
)
