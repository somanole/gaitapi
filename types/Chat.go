// Chat
package types

import "code.google.com/p/go-uuid/uuid"

type Chat struct{
	SenderId uuid.UUID
	ReceiverId uuid.UUID
	ReceiverUsername string
	IsChatActive bool
	IsChatBlockedBySender bool
	IsChatBlockedByReceiver bool
	LastMessage string
	Timestamp int64
}

type Chats []Chat

type ChatRequest struct {
	FirstUserId string
	SecondUserId string
}
