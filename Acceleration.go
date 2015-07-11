// Acceleration
package main

type Acceleration struct {
	Id int `json:uid`
	Timestamp int64 `json:timestamp`
	X float64 `json:x`
	Y float64 `json:y`
	Z float64 `json:z`
}

type Accelerations []Acceleration
