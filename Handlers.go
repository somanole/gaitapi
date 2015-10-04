// Handlers
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
	"io"
	"io/ioutil"
	"log"
	"github.com/gorilla/mux"
	"github.com/somanole/gaitapi/types"
	"github.com/somanole/gaitapi/accelerationrepo"
	"github.com/somanole/gaitapi/userrepo"
	"github.com/somanole/gaitapi/matchrepo"
	"github.com/somanole/gaitapi/messagerepo"
	"github.com/somanole/gaitapi/activityrepo"
	"github.com/somanole/gaitapi/services"
)

var accelerationRepo accelerationrepo.AccelerationRepo
var userRepo userrepo.UserRepo
var matchRepo matchrepo.MatchRepo
var messageRepo messagerepo.MessageRepo
var activityRepo activityrepo.ActivityRepo

func init() {
	accelerationRepo = accelerationrepo.New()
	userRepo = userrepo.New()
	matchRepo = matchrepo.New()
	messageRepo = messagerepo.New()
	activityRepo = activityrepo.New()
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Let's go for a walk.")
}

func CreateUserAcceleration(w http.ResponseWriter, r *http.Request) {
	var acceleration types.AccelerationRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	log.Printf(fmt.Sprintf("CreateUserAcceleration - Received userId: %v", userId))
	
	if userId != "" {
		 body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
	    if err != nil {
	        hasexploded = err.Error()
	    } else if err := r.Body.Close(); err != nil {
	        hasexploded = err.Error()
	    } else if err := json.Unmarshal(body, &acceleration); err != nil {
	        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	        w.WriteHeader(422) // unprocessable entity
			
	        if err := json.NewEncoder(w).Encode(err); err != nil {
	            hasexploded = err.Error()
	        }
	    } else {
			err := accelerationRepo.CreateAcceleration(userId, acceleration)
	    	
			if err != nil {
				if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotAcceptable)	
				} else if err.Error() == "not uuid" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    		w.WriteHeader(http.StatusBadRequest)	
				} else {
					hasexploded = err.Error()
				}		
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    	w.WriteHeader(http.StatusCreated)
			}   
		}
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("AccelerationCreate - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func ValidateAccessCode(w http.ResponseWriter, r *http.Request) {
    var accessCode services.AccessCode
	var hasexploded string
	hasexploded = ""
	
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
    if err != nil {
        hasexploded = err.Error()
    } else if err := r.Body.Close(); err != nil {
        hasexploded = err.Error()
    } else if err := json.Unmarshal(body, &accessCode); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
		
        if err := json.NewEncoder(w).Encode(err); err != nil {
            hasexploded = err.Error()
        }
    } else {
		response := services.ValidateAccessCode(accessCode)
	    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(http.StatusOK)
		
	    if err := json.NewEncoder(w).Encode(response); err != nil {
	        hasexploded = err.Error()
	    }
	} 
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("ValidateAccessCode - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func GetAccessCode(w http.ResponseWriter, r *http.Request) {
    response := services.GetAccessCode()
	
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf(fmt.Sprintf("GetAccessCode - Error: %v", err.Error()))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
    }
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserRequest
	var hasexploded string
	hasexploded = ""
	
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
    if err != nil {
        hasexploded = err.Error()
    } else if err := r.Body.Close(); err != nil {
        hasexploded = err.Error()
    } else if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            hasexploded = err.Error()
        }
    } else {
		response, err := userRepo.CreateUser(user)
	
		if err != nil {
			if err.Error() == "email already registered" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    	w.WriteHeader(http.StatusConflict)	
			} else {
				hasexploded = err.Error()
			}	
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    	w.WriteHeader(http.StatusCreated)
	    	if err := json.NewEncoder(w).Encode(response); err != nil {
	        	hasexploded = err.Error()
	    	}
		}   
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("CreateUser - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	if userId != "" {
		var response types.GetUserResponse
	
	    user, err := userRepo.GetUser(userId)
		
		if err != nil{
			if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotAcceptable)	
			} else if err.Error() == "not uuid" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusBadRequest)	
			} else {
				hasexploded = err.Error()
			}		
		} else {
			response.DeviceType = user.DeviceType
			response.IsAnonymous = user.IsAnonymous
			response.IsTest = user.IsTest
			response.Timestamp = user.Timestamp
			response.UserId = user.UserId
			response.Username = user.Username
			
		    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        hasexploded = err.Error()
		    }
		}	
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUser - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user types.UserUpdateRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]	
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
    if err != nil {
        hasexploded = err.Error()
    } else if err := r.Body.Close(); err != nil {
        hasexploded = err.Error()
    } else if err := json.Unmarshal(body, &user); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            hasexploded = err.Error()	
        }
    } else {
		response, err := userRepo.UpdateUser(userId, user)
	
		if err != nil {
			if err.Error() == "not found" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    	w.WriteHeader(http.StatusNotAcceptable)	
			} else if err.Error() == "not uuid" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    		w.WriteHeader(http.StatusBadRequest)	
			} else {
				hasexploded = err.Error()
			}
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
			
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        hasexploded = err.Error()
		    }
		}
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("UpdateUser - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	email := vars["email"]
	
	log.Printf(fmt.Sprintf("GetUserByEmail handler - email received: %v", email))
	
	if email != "" {
	    response, err := userRepo.GetUserByEmail(email)
		
		if err != nil{
			if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotAcceptable)	
			} else {
				hasexploded = err.Error()
			}		
		} else {
		    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
			
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        hasexploded = err.Error()
		    }
		}	
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserByEmail - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetUserExtraInfo(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	if userId != "" {
	    response, err := userRepo.GetUserExtraInfo(userId)
		
		if err != nil{
			if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotAcceptable)	
			} else if err.Error() == "not uuid" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusBadRequest)	
			} else {
				hasexploded = err.Error()
			}		
		} else {
		    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
			
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        hasexploded = err.Error()
		    }
		}	
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserExtraInfo - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetUserMatch(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	if userId != "" {
	    response, err := matchRepo.GetUserMatch(userId)
		
		if err != nil{
			if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNoContent)	
			} else if err.Error() == "not uuid" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusBadRequest)	
			} else {
				hasexploded = err.Error()
			}		
		} else {
		    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
			
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        hasexploded = err.Error()
		    }
		}	
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserMatch - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message types.MessageRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	receiverId := vars["receiverid"]
	
	if userId != "" && receiverId != "" {
		 body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
	    if err != nil {
	        hasexploded = err.Error()
	    } else if err := r.Body.Close(); err != nil {
	        hasexploded = err.Error()
	    } else if err := json.Unmarshal(body, &message); err != nil {
	        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	        w.WriteHeader(422) // unprocessable entity
	        if err := json.NewEncoder(w).Encode(err); err != nil {
	            hasexploded = err.Error()
	        }
	    } else {
			err := messageRepo.CreateMessage(userId, receiverId, message)
		
			if err != nil {
				if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotAcceptable)	
				} else if err.Error() == "not uuid" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    		w.WriteHeader(http.StatusBadRequest)	
				} else {
					hasexploded = err.Error()
				}		
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    	w.WriteHeader(http.StatusCreated)
			}   
		}
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("CreateMessage - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func GetUserMessagesByReceiverId(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	startdate := r.URL.Query().Get("startdate")
	log.Printf(fmt.Sprintf("GetUserMatch - startdate: %v", startdate))
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	receiverId := vars["receiverid"]
	
	if userId != "" && receiverId != "" {
	    response, err := messageRepo.GetUserMessagesByReceiverId(userId, receiverId, startdate)
		
		if err != nil{
			if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotAcceptable)	
			} else if err.Error() == "not uuid" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusBadRequest)	
			} else {
				hasexploded = err.Error()
			}		
		} else if response != nil {
		    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
			
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        hasexploded = err.Error()
		    }
		} else {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusNoContent)
		}	
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserMatch - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateUserActivity(w http.ResponseWriter, r *http.Request) {
	var activity types.ActivityRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	if userId != "" {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
	    if err != nil {
	        hasexploded = err.Error()
	    } else if err := r.Body.Close(); err != nil {
	        hasexploded = err.Error()
	    } else if err := json.Unmarshal(body, &activity); err != nil {
	        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	        w.WriteHeader(422) // unprocessable entity
	        if err := json.NewEncoder(w).Encode(err); err != nil {
	            hasexploded = err.Error()
	        }
	    } else {
			err := activityRepo.CreateUserActivity(userId, activity)
		
			if err != nil {
				if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNotAcceptable)	
				} else if err.Error() == "not uuid" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    		w.WriteHeader(http.StatusBadRequest)	
				} else {
					hasexploded = err.Error()
				}		
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    	w.WriteHeader(http.StatusCreated)
			}   
		}
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("CreateUserActivity - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func GetUserActivity(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	if userId != "" {
	    response, err := activityRepo.GetUserActivity(userId)
		
		if err != nil{
			if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    			w.WriteHeader(http.StatusNoContent)	
			} else if err.Error() == "not uuid" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusBadRequest)	
			} else {
				hasexploded = err.Error()
			}		
		} else {
		    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		    w.WriteHeader(http.StatusOK)
			
		    if err := json.NewEncoder(w).Encode(response); err != nil {
		        hasexploded = err.Error()
		    }
		}	
	} else {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
    	w.WriteHeader(http.StatusBadRequest)	
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserActivity - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
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

func HelpPagePOSTUser(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-user.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePUTUser(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/PUT-user.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETUserId(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-user-id.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETUserEmail(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-user-email.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETExtraInfo(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-extrainfo.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETMatch(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-match.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePOSTMessage(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-message.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETMessage(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-message.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePOSTActivity(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-activity.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETActivity(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-activity.html")
    fmt.Fprint(w, string(body))
}
