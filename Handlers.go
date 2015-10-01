// Handlers
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
	"io"
	"io/ioutil"
	"github.com/gorilla/mux"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/accelerationrepo"
	"github.com/somanole/gaitapi/userrepo"
	"github.com/somanole/gaitapi/services"
)

var accelerationRepo accelerationrepo.AccelerationRepo
var userRepo userrepo.UserRepo

func init() {
	accelerationRepo = accelerationrepo.New()
	userRepo = userrepo.New()
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Let's go for a walk.")
}

func AccelerationIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	
	accelerations := accelerationRepo.GetAllAccelerations()
    if err := json.NewEncoder(w).Encode(accelerations); err != nil {
        panic(err)
    }
}

func AccelerationsCount(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Method temporarily unavailable")
	
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    //w.WriteHeader(http.StatusOK)
	
	//count := accelerationRepo.GetAccelerationsCount()
    //if err := json.NewEncoder(w).Encode(count); err != nil {
		//panic(err)
    //    fmt.Fprintln(w, err)
    //}
}

func AccelerationCreate(w http.ResponseWriter, r *http.Request) {
	var acceleration types.Acceleration
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &acceleration); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    accelerationRepo.CreateAcceleration(acceleration)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
}

func ValidateAccessCode(w http.ResponseWriter, r *http.Request) {
    var accessCode services.AccessCode
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &accessCode); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    response := services.ValidateAccessCode(accessCode)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        panic(err)
    }
}

func GetAccessCode(w http.ResponseWriter, r *http.Request) {
    response := services.GetAccessCode()
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        panic(err)
    }
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    response := userRepo.CreateUser(user)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        panic(err)
    }
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	
	if userId != "" {
		var response types.GetUserResponse
	
	    user, err := userRepo.GetUser(userId)
		
		if err != nil{
			if err.Error() == "not found"{
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotFound)	
			} else if err.Error() == "not uuid"{
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusBadRequest)	
			} else{
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusInternalServerError)	
			}		
		} else{
			response.DeviceType = user.DeviceType
			response.IsAnonymous = user.IsAnonymous
			response.IsTest = user.IsTest
			response.Timestamp = user.Timestamp
			response.UserId = user.UserId
			response.Username = user.Username
			
		    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        panic(err)
		    }
		}	
	} else{
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	vars := mux.Vars(r)
	userId := vars["id"]
	
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
	
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
	
    if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    response, err := userRepo.UpdateUser(userId, user)
	
	if err != nil{
		if err.Error() == "not found"{
			w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    	w.WriteHeader(http.StatusNotFound)	
		} else if err.Error() == "not uuid"{
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusBadRequest)	
		} else{
			w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    	w.WriteHeader(http.StatusInternalServerError)	
		}
	} else{
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(http.StatusOK)
	    if err := json.NewEncoder(w).Encode(response); err != nil {
	        panic(err)
	    }
	}
}

func HelpPageIndex(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/index.html")
    fmt.Fprint(w, string(body))
}

func HelpPageCss(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/style.css")
	
	w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	
    fmt.Fprint(w, string(body))
}

func HelpPagePOSTAccesscodeValidate(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-accesscode-validate.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETAccesscode(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-accesscode.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePOSTAcceleration(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-acceleration.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETAccelerations(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-accelerations.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETAccelerationsCount(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-accelerations-count.html")
    fmt.Fprint(w, string(body))
}
