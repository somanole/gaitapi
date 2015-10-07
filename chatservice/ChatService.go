// ChatService
package chatservice

import (
	"errors"
	"time"
	"github.com/somanole/gaitapi/chatrepo"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
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
				if secondUsername, err = utilsservice.GetUserUsername(cr.FirstUserId); err == nil {
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
			err = DeleteChat(senderId, receiverId)
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

	return err
}
