// MessageRepo
package messagerepo

import (
	"github.com/somanole/gaitapi/types"
)

type (
	MessageRepo interface {
		GetUserMessagesByReceiverId(userId string, receiverId string, startdate string) (types.Messages, error)
		CreateMessage(m types.Message) error
	} 
	
	repoFactory func() MessageRepo
)

var (
	New repoFactory
)
