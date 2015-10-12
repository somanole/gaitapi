// AccelerationRepo
package accelerationrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	AccelerationRepo interface {
		CreateAcceleration(a types.Acceleration) error
	}
	
	repoFactory func() AccelerationRepo
)

var (
	New repoFactory
)
