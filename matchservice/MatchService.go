// MatchService
package matchservice

import (
	"time"
	"github.com/somanole/gaitapi/matchrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/chatservice"
	"github.com/somanole/gaitapi/notificationsservice"
	"github.com/somanole/gaitapi/activityservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
)

var matchRepo matchrepo.MatchRepo

func init() {
	matchRepo = matchrepo.New()
}

func CreateMatch(mr types.MatchRequest) error {
	var err error
	err = nil

	if err = utilsservice.CheckIfUserExists(mr.FirstUserId); err == nil {
		if err = utilsservice.CheckIfUserExists(mr.SecondUserId); err == nil {
			var firstUsername string
			var secondUsername string
			
			if firstUsername, err = utilsservice.GetUserUsername(mr.FirstUserId); err == nil {
				if secondUsername, err = utilsservice.GetUserUsername(mr.SecondUserId); err == nil {
					var m types.Match
			
					m.UserId = uuid.Parse(mr.FirstUserId)
					m.MatchedUserId = uuid.Parse(mr.SecondUserId)
					m.MatchedUsername = secondUsername
					m.Timestamp = int64(time.Now().UTC().Unix())
					
					//create first match
					if err = matchRepo.CreateMatch(m); err == nil {
						m.UserId = uuid.Parse(mr.SecondUserId)
						m.MatchedUserId = uuid.Parse(mr.FirstUserId)
						m.MatchedUsername = firstUsername
						m.Timestamp = int64(time.Now().UTC().Unix())
						
						//create second match
						if err = matchRepo.CreateMatch(m); err == nil {
							var cr types.ChatRequest
							
							//create chat
							cr.FirstUserId = mr.FirstUserId
							cr.SecondUserId = mr.SecondUserId
							err = chatservice.CreateChat(cr)
							
							//send push message to first user
							if lastActivity, errA := activityservice.GetUserActivity(mr.FirstUserId); errA == nil {
								notificationsservice.SendPushNotification(lastActivity.DeviceType, lastActivity.PushToken, "you've been matched!")
							}
							
							//send push message to second user
							if lastActivity, errA := activityservice.GetUserActivity(mr.SecondUserId); errA == nil {
								notificationsservice.SendPushNotification(lastActivity.DeviceType, lastActivity.PushToken, "you've been matched!")
							}
						}
					}
				}
			}
		}
	}
	
	return err
}

func GetUserMatch(userId string) (types.Match, error) {
	var response types.Match
	var err error
	err = nil

	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		response, err = matchRepo.GetUserMatch(userId)
	}
	
	return response, err
}
