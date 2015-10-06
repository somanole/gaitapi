// ChatRepo
package chatrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	ChatRepo interface {
		GetUserActiveChats(userId string) (types.Chats, error)
		BlockChat(senderId string, receiverId string) error
		UnblockChat(senderId string, receiverId string) error
		DeleteChat(senderId string, receiverId string) error
	} 
	
	repoFactory func() ChatRepo
)

var (
	New repoFactory
)
