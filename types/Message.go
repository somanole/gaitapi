// Message
package types

import "code.google.com/p/go-uuid/uuid"

type Message struct{
	MessageId uuid.UUID
	SenderId uuid.UUID
	ReceiverId uuid.UUID
	Text string
	IsRead bool
	Timestamp int64
}


