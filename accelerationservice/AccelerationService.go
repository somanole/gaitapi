// AccelerationService
package accelerationservice

import (
	"time"
	"github.com/somanole/gaitapi/accelerationrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
)

var accelerationRepo accelerationrepo.AccelerationRepo

func init() {
	accelerationRepo = accelerationrepo.New()
}

func CreateAccelerations(userId string, accelerations types.AccelerationsRequest) error {
	var err error
	err = nil
	
	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		for _,ar := range accelerations {
			var a types.Acceleration
			a.UserId = uuid.Parse(userId)
			a.X = ar.X
			a.Y = ar.Y
			a.Z = ar.Z
			
			if ar.Timestamp != 0 {
				a.Timestamp = int64(time.Unix(ar.Timestamp, 0).UTC().Unix())
			} else {
				a.Timestamp = int64(time.Now().UTC().Unix())
			}
			
			err = accelerationRepo.CreateAcceleration(a)
			
			if err != nil {
				break;
			}
		}
	}
	
	return err
}
