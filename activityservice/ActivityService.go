// ActivityService
package activityservice

import (
	"time"
	"github.com/somanole/gaitapi/activityrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
)

var activityRepo activityrepo.ActivityRepo

func init() {
	activityRepo = activityrepo.New()
}

func CreateUserActivity(userId string, ar types.ActivityRequest) error {
	var err error
	err = nil
	
	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		var a types.Activity
			
		a.DeviceId = ar.DeviceId
		a.DeviceType = ar.DeviceType
		a.IsLoggedIn = ar.IsLoggedIn
		a.UserId = uuid.Parse(userId)
		a.PushToken = ar.PushToken

		if ar.Timestamp != 0 {
			a.Timestamp = int64(time.Unix(ar.Timestamp, 0).UTC().Unix())
		} else {
			a.Timestamp = int64(time.Now().UTC().Unix())
		}
			
		err = activityRepo.CreateUserActivity(a)
	}
	
	return err
}

func GetUserActivity(userId string) (types.Activity, error) {
	var response types.Activity
	var err error
	err = nil
		
	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		response, err = activityRepo.GetUserActivity(userId)
	}
	
	return response, err
}
