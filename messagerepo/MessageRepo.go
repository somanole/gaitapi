// MessageRepo
package messagerepo

import (
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
)

type (
	MessageRepo interface {
		GetUserMessagesByReceiverId(userId string, receiverId string, startdate string) (types.Messages, error)
		CreateMessage(m types.Message) error
		DeleteMessage(senderId uuid.UUID, receiverId uuid.UUID, timestamp int64) error
		GetUserLastMessageByReceiverId(userId string, receiverId string) (types.Message, error)
	} 
	
	repoFactory func() MessageRepo
)

var (
	New repoFactory
)
