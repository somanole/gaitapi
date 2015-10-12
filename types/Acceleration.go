// Acceleration
package types

import "code.google.com/p/go-uuid/uuid"

type Acceleration struct {
	UserId uuid.UUID 
	Timestamp int64 
	X float64 
	Y float64 
	Z float64
}

type AccelerationRequest struct {
	Timestamp int64 
	X float64 
	Y float64 
	Z float64
}

type Accelerations []Acceleration

type AccelerationsRequest []AccelerationRequest

type AccelerationsCount struct {
	AccelerationsCount int64 
}
