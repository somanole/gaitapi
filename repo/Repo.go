// Repo
package repo

import (
	"github.com/somanole/gait/acceleration"
)

type (
	Repo interface {
		GetAcceleration(userId int64) acceleration.Acceleration
		GetAllAccelerations() acceleration.Accelerations
		CreateAcceleration(acceleration acceleration.Acceleration) acceleration.Acceleration
		DeleteAcceleration(userId int64) error
	}
	
	repoFactory func() Repo
)

var (
	New repoFactory
)
