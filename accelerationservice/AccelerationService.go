// AccelerationService
package accelerationservice

import (
	"time"
	"github.com/somanole/gaitapi/accelerationrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/userservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
	"log"
	"fmt"
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
		
		//increment walking progress for test users
		if err == nil {
			if user, errU := userservice.GetUser(userId); errU == nil {
				if user.IsTest {
					if userExtraInfo, errU := userservice.GetUserExtraInfo(userId); errU == nil {
						if userExtraInfo.WalkingProgress < 100 {
							var userExtraInfoRequest types.UserExtraInfoRequest
							
							if userExtraInfo.WalkingProgress >= 95 {
								userExtraInfoRequest.WalkingProgress = 100
							} else {
								userExtraInfoRequest.WalkingProgress = userExtraInfo.WalkingProgress + 3
							}
							
							if errU := userservice.UpdateUserExtraInfo(userId, userExtraInfoRequest); errU != nil {
								log.Printf(fmt.Sprintf("AccelerationService.CreateAccelerations() - Error: %v", errU.Error()))
							}
						}
					} else {
						log.Printf(fmt.Sprintf("AccelerationService.CreateAccelerations() - Error: %v", errU.Error()))
					}
				}
			} else {
				log.Printf(fmt.Sprintf("AccelerationService.CreateAccelerations() - Error: %v", errU.Error()))
			}
		}
	}
	
	return err
}
