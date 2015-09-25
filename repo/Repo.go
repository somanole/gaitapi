// Repo
package repo

import (
	"github.com/somanole/gaitapi/acceleration"
)

type (
	Repo interface {
		GetAcceleration(userId int64) acceleration.Acceleration
		GetAllAccelerations() acceleration.Accelerations
		GetAccelerationsCount() acceleration.AccelerationsCount
		CreateAcceleration(acceleration acceleration.Acceleration) acceleration.Acceleration
		DeleteAcceleration(userId int64) error
	}
	
	repoFactory func() Repo
)

var (
	New repoFactory
)
