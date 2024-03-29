// MessageService
package messageservice

import (
	"fmt"
	"strings"
	"time"
	"github.com/somanole/gaitapi/messagerepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/chatservice"
	"github.com/somanole/gaitapi/notificationsservice"
	"github.com/somanole/gaitapi/activityservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
)

var messageRepo messagerepo.MessageRepo

func init() {
	messageRepo = messagerepo.New()
}

func CreateMessage(userId string, receiverId string, mr types.MessageRequest) error {
	var err error
	err = nil

	if err = utilsservice.CheckIfMatchExists(userId, receiverId); err == nil {
		var m types.Message
		
		m.SenderId = uuid.Parse(userId)
		m.ReceiverId = uuid.Parse(receiverId)
		m.MessageId = uuid.NewRandom()
		m.IsRead = false
		m.Text = strings.Replace(mr.Text, "'", "''", -1)
		if mr.Timestamp != 0 {
			m.Timestamp = int64(time.Unix(mr.Timestamp, 0).UTC().Unix())
		} else {
			m.Timestamp = int64(time.Now().UTC().Unix())
		}
		
		if err = messageRepo.CreateMessage(m); err == nil {
			err = chatservice.UpdateLastMessageChat(userId, receiverId, m.Text)
			
			if lastActivity, errA := activityservice.GetUserActivity(receiverId); errA == nil {
				senderUsername, _ := utilsservice.GetUserUsername(userId)
				pushMessage := fmt.Sprintf("%s: %s", senderUsername, m.Text)
				notificationsservice.SendPushNotification(lastActivity.DeviceType, lastActivity.PushToken, pushMessage)
			}
		}
	}
	
	return err
}

func GetUserMessagesByReceiverId(userId string, receiverId string, startdate string) (types.Messages, error) {
	var response types.Messages
	var err error
	err = nil
	
	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		if err = utilsservice.CheckIfUserExists(userId); err == nil {	
			response, err = messageRepo.GetUserMessagesByReceiverId(userId, receiverId, startdate)
		}
	}
	
	return response, err
}

func GetUserLastMessageByReceiverId(userId string, receiverId string) (types.Message, error) {
	var response types.Message
	var err error
	err = nil
	
	if err = utilsservice.CheckIfUserExists(userId); err == nil {
		if err = utilsservice.CheckIfUserExists(userId); err == nil {	
			response, err = messageRepo.GetUserLastMessageByReceiverId(userId, receiverId)
		}
	}
	
	return response, err
}

func DeleteMessages(userId string, dmsr types.DeleteMessagesRequest) error {
	var err error
	err = nil
	
	if err = utilsservice.CheckIfUUID(userId); err == nil {		
		for _,dmr := range dmsr {		
			if err = messageRepo.DeleteMessage(dmr.SenderId, dmr.ReceiverId, dmr.Timestamp); err == nil {
				var lastMessage types.Message
				
				if lastMessage, err = GetUserLastMessageByReceiverId(dmr.SenderId.String(), dmr.ReceiverId.String()); err == nil {
					err = chatservice.UpdateLastMessageChat(dmr.SenderId.String(), dmr.ReceiverId.String(), lastMessage.Text)
				}
			}
		}
	}
	
	return err
}
