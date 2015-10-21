// UserService
package userservice

import (
	"fmt"
	"time"
	"github.com/somanole/gaitapi/userrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/activityservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
	"errors"
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

func UpdateUserExtraInfo(userId string, uer types.UserExtraInfoRequest) error {
	var err error
	err = nil
	
	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		if uer.WalkingProgress >= 0 && uer.WalkingProgress <= 100 {
			var ue types.UserExtraInfo
			
			ue.UserId = uuid.Parse(userId)
			ue.WalkingProgress = uer.WalkingProgress
			ue.Timestamp = int64(time.Now().UTC().Unix())
				
			err = userRepo.UpdateUserExtraInfo(ue)
		} else {
			err = errors.New("400")
		}
	}
	
	return err
}

func GetUserByEmail(email string) (types.UserByEmail, error) {
	return userRepo.GetUserByEmail(email)
}

func CreateUser(ur types.UserRequest) (types.CreateUserResponse, error) {
	var err error 
	var u types.User
	var userByEmail types.UserByEmail
	var response types.CreateUserResponse
	err = nil
	
	if (!ur.IsAnonymous) {
		userByEmail, err = GetUserByEmail(ur.Email)
	}
	
	if (ur.IsAnonymous || (!ur.IsAnonymous && err != nil && err.Error() == "not found")) {
		err = nil
		
		var username string		
		wordnikResponse, err := utilsservice.GenerateRandomUsername()
		if err != nil {
			username = "brown fox"
			err = nil
		} else if len(wordnikResponse) >= 2 {
			username = fmt.Sprintf("%v %v", wordnikResponse[0].Word, wordnikResponse[1].Word)
		} else {
			username = "brown fox"
		}
		
		u.UserId = uuid.NewRandom()
		u.Username = username
		u.Timestamp = int64(time.Now().UTC().Unix())
		u.DeviceType = ur.DeviceType
		u.Email = ur.Email
		u.FacebookAccessToken = ur.FacebookAccessToken
		u.GenderPreference = ur.GenderPreference
		u.GoogleAccessToken = ur.GoogleAccessToken
		u.IsAnonymous = ur.IsAnonymous
		u.IsTest = false
		u.Password = ur.Password
		u.PushToken = ur.PushToken
		u.TwitterAccessToken = ur.TwitterAccessToken
				
		response, err = userRepo.CreateUser(u)
	} else if userByEmail.Email != "" {
		err = errors.New("email already registered")
	}
	
	return response, err
}
