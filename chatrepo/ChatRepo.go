// ChatRepo
package chatrepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	ChatRepo interface {
		CreateChat(chat types.Chat) error
		GetUserActiveChats(userId string) (types.Chats, error)
		BlockChat(senderId string, receiverId string) error
		UnblockChat(senderId string, receiverId string) error
		DeleteChat(senderId string, receiverId string) error
		UpdateLastMessageChat(senderId string, receiverId string, message string) error	
		CreateMatch(m types.Match) error
		GetUserMatch(userId string) (types.Match, error)
		GetUserPerfectNumber(userId string) (types.PerfectNumber, error)
		GetUserPerfectMatch(rpf types.PerfectNumber) (types.PerfectNumber, error) 
	} 
	
	repoFactory func() ChatRepo
)

var (
	New repoFactory
)
