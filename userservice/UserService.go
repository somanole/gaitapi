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

func GetUser(userId string) (types.User, error) {
	var err error
	var response types.User 
	err = nil
	
	if err = utilsservice.CheckIfUUID(userId); err == nil {
		response, err = userRepo.GetUser(userId)
	}
	
	return response, err
}

func GetUserExtraInfo(userId string) (types.UserExtraInfo, error) {
	var err error
	var response types.UserExtraInfo 
	err = nil
	
	if err = utilsservice.CheckIfUUID(userId); err == nil {
		response, err = userRepo.GetUserExtraInfo(userId)
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

func ReportUser(userId string, urr types.UserReportRequest) error {
	var err error
	err = nil
	
	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		var ur types.UserReport
		
		ur.ReportedUserId = uuid.Parse(userId)
		ur.ReporterUserId = urr.ReporterUserId
		ur.Reason = urr.Reason
		ur.Comment = urr.Comment
		ur.Timestamp = int64(time.Now().UTC().Unix())
		
		err = userRepo.ReportUser(ur)
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

func UpdateUser(userId string, u types.UserUpdateRequest) (types.CreateUserResponse, error) {
	var err error
	var response types.CreateUserResponse 
	err = nil
	
	if err = utilsservice.CheckIfUUID(userId); err == nil {
		var user types.User
		
		if user, err = GetUser(userId); err == nil {
			user.Timestamp = int64(time.Now().UTC().Unix())
			user.IsTest = u.IsTest
			
			if u.FacebookAccessToken != "" { 
				user.FacebookAccessToken = u.FacebookAccessToken
			}
			if u.DeviceType != "" {
				user.DeviceType = u.DeviceType
			}
			if u.GenderPreference != "" {
				user.GenderPreference = u.GenderPreference
			}
			if u.GoogleAccessToken != "" {
				user.GoogleAccessToken = u.GoogleAccessToken
			}
			if u.Password != "" {
				user.Password = u.Password
			}
			if u.PushToken != "" {
				user.PushToken = u.PushToken
			}
			if u.TwitterAccessToken != "" {
				user.TwitterAccessToken = u.TwitterAccessToken
			}
			
			response, err = userRepo.UpdateUser(user)
		}
	}
	
	return response, err
}
