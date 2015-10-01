// Activity
package types

import "code.google.com/p/go-uuid/uuid"

type Activity struct{
	UserId uuid.UUID
	DeviceType string
	DeviceId string
	IsLoggedIn bool
	pushToken string
	Timestamp int64
}
