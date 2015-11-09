// ChatService
package chatservice

import (
	"errors"
	"time"
	"github.com/somanole/gaitapi/chatrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/activityservice"
	"github.com/somanole/gaitapi/notificationsservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
	"log"
	"fmt"
)

var chatRepo chatrepo.ChatRepo

func init() {
	chatRepo = chatrepo.New()
}

func CreateChat(cr types.ChatRequest) error {
	var err error
	err = nil

	if err = utilsservice.CheckIfUserExists(cr.FirstUserId); err == nil {
		if err = utilsservice.CheckIfUserExists(cr.SecondUserId); err == nil {
			var firstUsername string
			var secondUsername string
				
			if firstUsername, err = utilsservice.GetUserUsername(cr.FirstUserId); err == nil {
				if secondUsername, err = utilsservice.GetUserUsername(cr.SecondUserId); err == nil {
					var c types.Chat
						
					c.SenderId = uuid.Parse(cr.FirstUserId)
					c.ReceiverId = uuid.Parse(cr.SecondUserId)
					c.ReceiverUsername = secondUsername
					c.IsChatActive = true
					c.IsChatBlockedBySender = false
					c.IsChatBlockedByReceiver = false
					c.LastMessage = ""
					c.Timestamp = int64(time.Now().UTC().Unix())
						
					if err = chatRepo.CreateChat(c); err == nil {
						c.SenderId = uuid.Parse(cr.SecondUserId)
						c.ReceiverId = uuid.Parse(cr.FirstUserId)
						c.ReceiverUsername = firstUsername
						c.IsChatActive = true
						c.IsChatBlockedBySender = false
						c.IsChatBlockedByReceiver = false
						c.LastMessage = ""
						c.Timestamp = int64(time.Now().UTC().Unix())
							
						err = chatRepo.CreateChat(c)
					}
				}
			}
		}
	}
	
	return err
}

func GetUserActiveChats(userId string) (types.Chats, error) {
	var response types.Chats
	var err error
	err = nil

	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		CreateUserPerfectMatch(userId)
		response, err = chatRepo.GetUserActiveChats(userId)
	}
	
	return response, err
}

func UpdateChat (senderId string, receiverId string, action string) error {
	var err error
	err = errors.New("405")
	
	switch action {
		case "block" :
			err = BlockChat(senderId, receiverId)
		case "unblock" :
			err = UnblockChat(senderId, receiverId)
		case "delete" :
			err = BlockChat(senderId, receiverId)
			if err == nil {
				err = DeleteChat(senderId, receiverId)
			}
	}
	
	return err
}

func BlockChat(senderId string, receiverId string) error {
	var err error
	err = nil

	if err = utilsservice.CheckIfMatchExists(senderId, receiverId); err == nil {
		err = chatRepo.BlockChat(senderId, receiverId)
	}
	
	return err
}

func UnblockChat(senderId string, receiverId string) error {
	var err error
	err = nil

	if err = utilsservice.CheckIfMatchExists(senderId, receiverId); err == nil {
		err = chatRepo.UnblockChat(senderId, receiverId)
	}
	
	return err
}

func DeleteChat(senderId string, receiverId string) error {
	var err error
	err = nil

	if err = utilsservice.CheckIfMatchExists(senderId, receiverId); err == nil {
		err = chatRepo.DeleteChat(senderId, receiverId)
	}
	
	return err
}

func UpdateLastMessageChat(senderId string, receiverId string, message string) error {
	var err error
	err = nil
	
	err = chatRepo.UpdateLastMessageChat(senderId, receiverId, message)
	err = chatRepo.UpdateLastMessageChat(receiverId, senderId, message)

	return err
}

func CreateUserPerfectMatch(userId string) error {
	var err error
	err = nil
		
	if _, err = chatRepo.GetUserMatch(userId); err != nil {
		var perfectNumber types.PerfectNumber
		if perfectNumber, err = chatRepo.GetUserPerfectNumber(userId); err == nil {
			var perfectMatch types.PerfectNumber
				
			if perfectMatch, err = chatRepo.GetUserPerfectMatch(perfectNumber); err == nil {
				if err = utilsservice.CheckIfUUID(perfectMatch.UserId.String()); err == nil {
					var mr types.MatchRequest
					mr.FirstUserId = userId
					mr.SecondUserId = perfectMatch.UserId.String()
						
					err = CreateMatch(mr)
				} else {
					log.Printf(fmt.Sprintf("ChatService.CreateUserPerfectMatch() - Error: %v", err.Error()))
				}		
			}		
		}
	}
	
	return err
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
					if err = chatRepo.CreateMatch(m); err == nil {
						m.UserId = uuid.Parse(mr.SecondUserId)
						m.MatchedUserId = uuid.Parse(mr.FirstUserId)
						m.MatchedUsername = firstUsername
						m.Timestamp = int64(time.Now().UTC().Unix())
						
						//create second match
						if err = chatRepo.CreateMatch(m); err == nil {
							var cr types.ChatRequest
							
							//create chat
							cr.FirstUserId = mr.FirstUserId
							cr.SecondUserId = mr.SecondUserId
							err = CreateChat(cr)
							
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
		response, err = chatRepo.GetUserMatch(userId)
	}
	
	return response, err
}
