// +build test

// MockedAccelerationRepo
package accelerationrepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
	"code.google.com/p/go-uuid/uuid"
)

type MockedAccelerationRepo struct {}

func NewMockedAccelerationRepo() AccelerationRepo {
	return &MockedAccelerationRepo{}
}

func init() {
	New = NewMockedAccelerationRepo
}

var currentId int64
var accelerations types.Accelerations

func (repo *MockedAccelerationRepo) CreateAcceleration(a types.Acceleration) error {
    var err error
	err = nil
	
    a.UserId = uuid.NewRandom()
    accelerations = append(accelerations, a)
    return err
}
