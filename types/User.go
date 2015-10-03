// User
package types

import "code.google.com/p/go-uuid/uuid"

type User struct{
	UserId uuid.UUID
	Username string
	FacebookAccessToken string
	TwitterAccessToken string
	GoogleAccessToken string
	PushTokeniOS string
	PushTokenAndroid string
	DeviceType string
	Email string
	Password string
	IsTest bool
	IsAnonymous bool
	GenderPreference string
	Timestamp int64
}

type UserRequest struct{
	FacebookAccessToken string
	TwitterAccessToken string
	GoogleAccessToken string
	PushTokeniOS string
	PushTokenAndroid string
	DeviceType string
	Email string
	Password string
	IsAnonymous bool
	GenderPreference string
}

type UserUpdateRequest struct{
	FacebookAccessToken string
	TwitterAccessToken string
	GoogleAccessToken string
	PushTokeniOS string
	PushTokenAndroid string
	DeviceType string
	Password string
	IsAnonymous bool
	GenderPreference string
}

type UserByEmail struct{
	Email string
	UserId uuid.UUID
}

type UserExtraInfo struct{
	UserId uuid.UUID
	WalkingProgress int
	Timestamp int64
}

type CreateUserResponse struct{
	UserId uuid.UUID
	Username string
	Timestamp int64
}

type GetUserResponse struct{
	UserId uuid.UUID
	Username string
	DeviceType string
	IsTest bool
	IsAnonymous bool
	Timestamp int64
}
