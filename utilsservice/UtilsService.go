// UtilsService
package utilsservice

import (
	"errors"
	"fmt"
	"log"
	"github.com/somanole/gaitapi/utilsrepo"
	"code.google.com/p/go-uuid/uuid"
	"github.com/somanole/gaitapi/types"
	"time"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
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

func RegisterInterest(ir types.InterestRequest) error {
	var i types.Interest
	var err error
	err = nil
	
	log.Printf(fmt.Sprintf("UtilsService.RegisterInterest() - Received email: %v", ir.Email))
	
	if ir.Email != "" {
		i.Timestamp = int64(time.Now().UTC().Unix())
		i.Email = ir.Email
			
		err = utilsRepo.RegisterInterest(i)		
	} else {
		log.Printf(fmt.Sprintf("UtilsService.RegisterInterest() - Received email is blank: %v", ir.Email))
		err = errors.New("blank credentials")
	}
	
	return err
}

func GenerateRandomUsername() (types.WordnikResponse, error) {
	log.Printf("UtilsService.GenerateRandomUsername() - ENTERED!!")
	
	var errorReturn error
	var wordnikResponse types.WordnikResponse
	errorReturn = nil
	
	r, err := http.Get("http://api.wordnik.com/v4/words.json/randomWords?hasDictionaryDef=true&minCorpusCount=0&minLength=3&maxLength=15&limit=2&api_key=a2a73e7b926c924fad7001ca3111acd55af2ffabf50eb4ae5")
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
	if err != nil {
	    errorReturn = err
        log.Printf(fmt.Sprintf("UtilsService.GenerateRandomUsername() - Error: %s", err))
	} else if err := r.Body.Close(); err != nil {
	    errorReturn = err
        log.Printf(fmt.Sprintf("UtilsService.GenerateRandomUsername() - Error: %s", err))
	} else if err := json.Unmarshal(body, &wordnikResponse); err != nil {
	    errorReturn = err
        log.Printf(fmt.Sprintf("UtilsService.GenerateRandomUsername() - Error: %s", err))
	} else {
		log.Printf(fmt.Sprintf("UtilsService.GenerateRandomUsername() - RESPONSE: +%v", wordnikResponse))
	}
	
	
	log.Printf("UtilsService.GenerateRandomUsername() - EXIT!!")
	
	return wordnikResponse, errorReturn
}
