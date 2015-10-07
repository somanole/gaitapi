// Match
package types

import "code.google.com/p/go-uuid/uuid"

type Match struct {
	UserId uuid.UUID
	MatchedUserId uuid.UUID
	MatchedUsername string
	Timestamp int64
}

type Matches []Match

type MatchGroup struct {
	Groupname string
	UserId uuid.UUID
	hashPrefix int
	Timestamp int64
}

type MatchRequest struct {
	FirstUserId string
	SecondUserId string
}
