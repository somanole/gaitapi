// Acceleration
package acceleration

type Acceleration struct {
	UserId int64 `json:userid`
	Timestamp int64 `json:timestamp`
	X float64 `json:x`
	Y float64 `json:y`
	Z float64 `json:z`
}

type Accelerations []Acceleration

type AccelerationsCount struct {
	AccelerationsCount int64 
}
