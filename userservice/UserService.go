// UserService
package userservice

import (
	"time"
	"github.com/somanole/gaitapi/userrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/activityservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
)

var userRepo userrepo.UserRepo

func init() {
	userRepo = userrepo.New()
}

func Login(l types.LoginRequest) (types.LoginResponse, error) {
	var userId string
	var err error
	var response types.LoginResponse
	err = nil
	
	if userId, err = utilsservice.CheckLoginCredentials(l.Email, l.Password); err == nil {
		var a types.ActivityRequest
		
		a.DeviceId = l.DeviceId
		a.DeviceType = l.DeviceType
		a.IsLoggedIn = l.IsLoggedIn
		a.PushToken = l.PushToken
		a.Timestamp = l.Timestamp
		
		err = activityservice.CreateUserActivity(userId, a)
		
		response.UserId = uuid.Parse(userId)
		response.Username, err = utilsservice.GetUserUsername(userId)
		response.Timestamp = int64(time.Now().UTC().Unix())
	}
	
	return response, err
}
