// MessageRepo
package messagerepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	MessageRepo interface {
		GetUserMessagesByReceiverId(userId string, receiverId string, startdate string) (types.Messages, error)
		CreateMessage(userId string, receiverId string, m types.MessageRequest) error
	} 
	
	repoFactory func() MessageRepo
)

var (
	New repoFactory
)
