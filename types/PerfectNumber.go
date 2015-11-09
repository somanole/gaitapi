// PerfectNumber
package types

import "code.google.com/p/go-uuid/uuid"

type PerfectNumber struct {
	UserId uuid.UUID
	Timestamp int64
	MeanX float64
	MeanY float64
	MeanZ float64
	VarianceX float64
	VarianceY float64
	VarianceZ float64
	AvgAbsDiffX float64
	AvgAbsDiffY float64
	AvgAbsDiffZ float64
	Resultant float64
	AvgTimePeakDiffY float64
	PerfectNumber float64
}

type PerfectNumbers []PerfectNumber
