// UtilsService
package utilsservice

import (
	"errors"
	"fmt"
	"log"
	"github.com/somanole/gaitapi/utilsrepo"
	"code.google.com/p/go-uuid/uuid"
)

var utilsRepo utilsrepo.UtilsRepo

func init() {
	utilsRepo = utilsrepo.New()
}

func CheckIfUUID(id string) error {
	var err error
	err = nil
	
	if uuid.Parse(id) == nil {
		log.Printf(fmt.Sprintf("UtilsService.CheckIfUUID() - Received Id: %v is not UUID", id))
		err = errors.New("not uuid")
	}
	
	return err
}

func CheckUserPassword(userId string, password string) error {
	var err error
	err = nil
	
	log.Printf(fmt.Sprintf("UtilsService.CheckUserPassword() - Received userId: %v", userId))
	log.Printf(fmt.Sprintf("UtilsService.CheckUserPassword() - Received password: %v", password))
	
	if err = CheckIfUUID(userId); err == nil {
		err = utilsRepo.CheckUserPassword(userId, password)
	} 
	
	return err
}

func CheckIfUserExists(userId string) error {
	var err error
	err = nil
	
	log.Printf(fmt.Sprintf("UtilsService.CheckIfUserExists() - Received userId: %v", userId))
	
	if err = CheckIfUUID(userId); err == nil {
		err = utilsRepo.CheckIfUserExists(userId)
	}
	
	return err
}

func CheckIfMatchExists(firstUserId string, secondUserId string) error {
	var err error
	err = nil
	
	log.Printf(fmt.Sprintf("UtilsService.CheckIfMatchExists() - Received firstUserId: %v", firstUserId))
	log.Printf(fmt.Sprintf("UtilsService.CheckIfMatchExists() - Received secondUserId: %v", secondUserId))
	
	if err = CheckIfUUID(firstUserId); err == nil {
		if err = CheckIfUUID(secondUserId); err == nil {
			err = utilsRepo.CheckIfMatchExists(firstUserId, secondUserId)
		}
	}
	
	return err
}

func GetUserUsername(userId string) (string, error) {
	var username string
	var err error
	err = nil
	
	log.Printf(fmt.Sprintf("UtilsService.GetUserUsername() - Received userId: %v", userId))
	
	if err = CheckIfUUID(userId); err == nil {
		username, err = utilsRepo.GetUserUsername(userId)
	}
	
	return username, err
}

func CheckLoginCredentials(email string, password string) (string, error) {
	var userId string
	var err error
	err = nil
	
	log.Printf(fmt.Sprintf("UtilsService.CheckLoginCredentials() - Received email: %v", email))
	log.Printf(fmt.Sprintf("UtilsService.CheckLoginCredentials() - Received password: %v", password))
	
	if email != "" && password != "" {
		userId, err = utilsRepo.CheckLoginCredentials(email, password)		
	} else {
		log.Printf(fmt.Sprintf("UtilsService.CheckLoginCredentials() - Received email or password is blank: %v, %v", email, password))
		err = errors.New("blank credentials")
	}
	
	return userId, err
}
