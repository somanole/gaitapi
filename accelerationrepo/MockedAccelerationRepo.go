// +build test

// MockedAccelerationRepo
package accelerationrepo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/types"
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

func (repo *MockedAccelerationRepo) GetAcceleration(userId int64) types.Acceleration {
    for _, a := range accelerations {
        if a.UserId == userId {
            return a
        }
    }
    // return empty Acceleration if not found
    return types.Acceleration{}
}

func (repo *MockedAccelerationRepo) GetAllAccelerations() types.Accelerations {
	log.Println("Mocked - trying to get all accelerations")
	
	return accelerations
}

func (repo *MockedAccelerationRepo) GetAccelerationsCount() types.AccelerationsCount {
	count := types.AccelerationsCount{len(accelerations)}
	return count
}

func (repo *MockedAccelerationRepo) CreateAcceleration(a types.Acceleration) types.Acceleration {
    currentId += 1
    a.UserId = currentId
    accelerations = append(accelerations, a)
    return a
}

func (repo *MockedAccelerationRepo) DeleteAcceleration(userId int64) error {
    for i, a := range accelerations {
        if a.UserId == userId {
            accelerations = append(accelerations[:i], accelerations[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Acceleration with userid of %d to delete", userId)
}
