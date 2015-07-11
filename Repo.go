// Repo
package main

import "fmt"

var currentId int

var accelerations Accelerations

// Give us some seed data
func init() {
    RepoCreateAcceleration(Acceleration{Id: 100})
    RepoCreateAcceleration(Acceleration{Id: 101})
}

func RepoFindAcceleration(id int) Acceleration {
    for _, a := range accelerations {
        if a.Id == id {
            return a
        }
    }
    // return empty Acceleration if not found
    return Acceleration{}
}

func RepoCreateAcceleration(a Acceleration) Acceleration {
    currentId += 1
    a.Id = currentId
    accelerations = append(accelerations, a)
    return a
}

func RepoDestroyAcceleration(id int) error {
    for i, a := range accelerations {
        if a.Id == id {
            accelerations = append(accelerations[:i], accelerations[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("Could not find Acceleration with id of %d to delete", id)
}
