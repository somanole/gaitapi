// Acceleration
package types

type Acceleration struct {
	UserId int64 
	Timestamp int64 
	X float64 
	Y float64 
	Z float64
}

type Accelerations []Acceleration

type AccelerationsCount struct {
	AccelerationsCount int64 
}
