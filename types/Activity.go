// Activity
package types

import "code.google.com/p/go-uuid/uuid"

type Activity struct{
	UserId uuid.UUID
	DeviceType string
	DeviceId string
	IsLoggedIn bool
	PushToken string
	Timestamp int64
}

type ActivityRequest struct{
	DeviceType string
	DeviceId string
	IsLoggedIn bool
	PushToken string
	Timestamp int64
}
