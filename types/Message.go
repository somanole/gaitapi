// Message
package types

import "code.google.com/p/go-uuid/uuid"

type Message struct {
	MessageId uuid.UUID
	SenderId uuid.UUID
	ReceiverId uuid.UUID
	Text string
	IsRead bool
	Timestamp int64
}

type Messages []Message

type MessageRequest struct {
	Text string
	Timestamp int64
}

type DeleteMessageRequest struct {
	SenderId uuid.UUID
	ReceiverId uuid.UUID
	Timestamp int64
}

type DeleteMessagesRequest []DeleteMessageRequest
