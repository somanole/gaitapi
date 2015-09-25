// +build test

// Repo
package repo

import (
	"fmt"
	"log"
	"github.com/somanole/gaitapi/acceleration"
)

type MockedRepo struct {}

func NewMockedRepo() Repo {
	return &MockedRepo{}
}

func init() {
	New = NewMockedRepo
}

var currentId int64
var accelerations acceleration.Accelerations

func (repo *MockedRepo) GetAcceleration(userId int64) acceleration.Acceleration {
    for _, a := range accelerations {
        if a.UserId == userId {
            return a
        }
    }
    // return empty Acceleration if not found
    return acceleration.Acceleration{}
}

func (repo *MockedRepo) GetAllAccelerations() acceleration.Accelerations {
	log.Println("Mocked - trying to get all accelerations")
	
	return accelerations
}

func (repo *MockedRepo) GetAccelerationsCount() acceleration.AccelerationsCount {
	count := acceleration.AccelerationsCount{len(accelerations)}
	return count
}

func (repo *MockedRepo) CreateAcceleration(a acceleration.Acceleration) acceleration.Acceleration {
    currentId += 1
    a.UserId = currentId
    accelerations = append(accelerations, a)
    return a
}

func (repo *MockedRepo) DeleteAcceleration(userId int64) error {
    for i, a := range accelerations {
        if a.UserId == userId {
            accelerations = append(accelerations[:i], accelerations[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Acceleration with userid of %d to delete", userId)
}
