// AccelerationRepo
package accelerationrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	AccelerationRepo interface {
		GetAcceleration(userId int64) types.Acceleration
		GetAllAccelerations() types.Accelerations
		GetAccelerationsCount() types.AccelerationsCount
		CreateAcceleration(acceleration types.Acceleration) types.Acceleration
		DeleteAcceleration(userId int64) error
	}
	
	repoFactory func() AccelerationRepo
)

var (
	New repoFactory
)
