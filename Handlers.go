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
	"github.com/somanole/gaitapi/matchservice"
	"github.com/somanole/gaitapi/messageservice"
	"github.com/somanole/gaitapi/accesscodeservice"
	"github.com/somanole/gaitapi/activityservice"
	"github.com/somanole/gaitapi/accelerationservice"
	"github.com/somanole/gaitapi/userservice"
	"github.com/somanole/gaitapi/chatservice"
	"github.com/somanole/gaitapi/utilsservice"
	"github.com/somanole/gaitapi/kafkaservice"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Let's go for a walk.")
}

func CreateUserAcceleration(w http.ResponseWriter, r *http.Request) {
	var accelerations types.AccelerationsRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
			if userId != "" {
			body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		
		    if err != nil {
		        hasexploded = err.Error()
		    } else if err := r.Body.Close(); err != nil {
		        hasexploded = err.Error()
		    } else if err := json.Unmarshal(body, &accelerations); err != nil {
		        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		        w.WriteHeader(422) // unprocessable entity
				
		        if err := json.NewEncoder(w).Encode(err); err != nil {
		            hasexploded = err.Error()
		        }
		    } else {
				err := accelerationservice.CreateAccelerations(userId, accelerations)
		    	
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}

	if hasexploded != "" {
		log.Printf(fmt.Sprintf("AccelerationCreate - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func GetAcceleration(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" {
		    response, err := accelerationservice.GetAccelerations(userId)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("Handlers.GetAccelerations() - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func ValidateAccessCode(w http.ResponseWriter, r *http.Request) {
    var accessCode accesscodeservice.AccessCode
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
		response := accesscodeservice.ValidateAccessCode(accessCode)
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
    response := accesscodeservice.GetAccessCode()
	
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
		response, err := userservice.CreateUser(user)
	
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

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginRequest types.LoginRequest
	var hasexploded string
	hasexploded = ""
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
	if err != nil {
	    hasexploded = err.Error()
	} else if err := r.Body.Close(); err != nil {
	    hasexploded = err.Error()
	} else if err := json.Unmarshal(body, &loginRequest); err != nil {
	    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(422) // unprocessable entity
	    if err := json.NewEncoder(w).Encode(err); err != nil {
	        hasexploded = err.Error()
	    }
	} else {
		response, err := userservice.Login(loginRequest)
		
		if err != nil {
			if err.Error() == "not found" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    		w.WriteHeader(http.StatusNotAcceptable)	
			} else if err.Error() == "blank credentials" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    	w.WriteHeader(http.StatusBadRequest)	
			} else if err.Error() == "401" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    	w.WriteHeader(http.StatusUnauthorized)	
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
		log.Printf(fmt.Sprintf("CreateUserActivity - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" {
			var response types.GetUserResponse
		
		    user, err := userservice.GetUser(userId)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
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
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
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
			response, err := userservice.UpdateUser(userId, user)
		
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
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
	    response, err := userservice.GetUserByEmail(email)
		
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
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" {
		    response, err := userservice.GetUserExtraInfo(userId)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}

	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserExtraInfo - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateUserExtraInfo(w http.ResponseWriter, r *http.Request) {
	var uer types.UserExtraInfoRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]	
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
    	if err != nil {
        	hasexploded = err.Error()
	    } else if err := r.Body.Close(); err != nil {
	        hasexploded = err.Error()
	    } else if err := json.Unmarshal(body, &uer); err != nil {
	        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	        w.WriteHeader(422) // unprocessable entity
	        if err := json.NewEncoder(w).Encode(err); err != nil {
	            hasexploded = err.Error()	
	        }
	    } else {
			err := userservice.UpdateUserExtraInfo(userId, uer)
		
			if err != nil {
				if err.Error() == "not found" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
			    	w.WriteHeader(http.StatusNotAcceptable)	
				} else if err.Error() == "not uuid" || err.Error() == "400" {
						w.Header().Set("Content-Type", "text/css; charset=UTF-8")
			    		w.WriteHeader(http.StatusBadRequest)	
				} else {
					hasexploded = err.Error()
				}
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			    w.WriteHeader(http.StatusOK)
			}
		}
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}

	if hasexploded != "" {
		log.Printf(fmt.Sprintf("Handlers.UpdateUserExtraInfo() - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func ReportUser(w http.ResponseWriter, r *http.Request) {
	var urr types.UserReportRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]	
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
    if err != nil {
        hasexploded = err.Error()
	} else if err := r.Body.Close(); err != nil {
	    hasexploded = err.Error()
	} else if err := json.Unmarshal(body, &urr); err != nil {
	    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(422) // unprocessable entity
	    if err := json.NewEncoder(w).Encode(err); err != nil {
	        hasexploded = err.Error()	
	    }
	} else {
		err := userservice.ReportUser(userId, urr)
		
		if err != nil {
			if err.Error() == "not found" {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
			    w.WriteHeader(http.StatusNotAcceptable)	
			} else if err.Error() == "not uuid" || err.Error() == "400" {
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

	if hasexploded != "" {
		log.Printf(fmt.Sprintf("Handlers.UpdateUserExtraInfo() - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateMatch(w http.ResponseWriter, r *http.Request) {
	var match types.MatchRequest
	var hasexploded string
	hasexploded = ""

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
	if err != nil {
	    hasexploded = err.Error()
	} else if err := r.Body.Close(); err != nil {
	    hasexploded = err.Error()
	} else if err := json.Unmarshal(body, &match); err != nil {
	    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(422) // unprocessable entity
	    if err := json.NewEncoder(w).Encode(err); err != nil {
	        hasexploded = err.Error()
	    }
	} else {
		err := matchservice.CreateMatch(match)
		
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

	if hasexploded != "" {
		log.Printf(fmt.Sprintf("CreateUserActivity - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)	
	}
}

func GetUserMatch(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" {
		    response, err := matchservice.GetUserMatch(userId)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserMatch - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetUserChats(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		response, err := chatservice.GetUserActiveChats(userId)
		
		if err != nil {
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
			if response != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusOK)
					
				if err := json.NewEncoder(w).Encode(response); err != nil {
				    hasexploded = err.Error()
				}
			} else {
				w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    	w.WriteHeader(http.StatusNoContent)	
			}
		}	
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}

	if hasexploded != "" {
		log.Printf(fmt.Sprintf("Handlers.GetUserChats() - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateChat(w http.ResponseWriter, r *http.Request) {
	var hasexploded string
	hasexploded = ""
	
	action := r.URL.Query().Get("action")
	log.Printf(fmt.Sprintf("Handlers.UpdateChat() - action: %v", action))
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	receiverId := vars["receiverid"]
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" && receiverId != "" && action != "" {
		    err := chatservice.UpdateChat(userId, receiverId, action)
			
			if err != nil{
				if err.Error() == "not found" {
						w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    			w.WriteHeader(http.StatusNotAcceptable)	
				} else if err.Error() == "not uuid" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    		w.WriteHeader(http.StatusBadRequest)	
				} else if err.Error() == "405" {
					w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    		w.WriteHeader(http.StatusMethodNotAllowed)	
				} else {
					hasexploded = err.Error()
				}		
			} else {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			    w.WriteHeader(http.StatusOK)
			}	
		} else {
			w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    	w.WriteHeader(http.StatusBadRequest)	
		}
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("Handlers.GetUserMatch() - HasExploded! - Error: %v", hasexploded))
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
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
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
				err := messageservice.CreateMessage(userId, receiverId, message)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
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
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" && receiverId != "" {
		    response, err := messageservice.GetUserMessagesByReceiverId(userId, receiverId, startdate)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserMatch - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	var dmsr types.DeleteMessagesRequest
	var hasexploded string
	hasexploded = ""
	
	vars := mux.Vars(r)
	userId := vars["userid"]
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" {
			body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		
		    if err != nil {
		        hasexploded = err.Error()
		    } else if err := r.Body.Close(); err != nil {
		        hasexploded = err.Error()
		    } else if err := json.Unmarshal(body, &dmsr); err != nil {
		        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		        w.WriteHeader(422) // unprocessable entity
		        if err := json.NewEncoder(w).Encode(err); err != nil {
		            hasexploded = err.Error()
		        }
		    } else {
				err := messageservice.DeleteMessages(userId, dmsr)
			
				if err != nil {
					if err.Error() == "not found" {
						w.Header().Set("Content-Type", "text/css; charset=UTF-8")
		    			w.WriteHeader(http.StatusNotAcceptable)	
					} else if err.Error() == "not uuid" {
						w.Header().Set("Content-Type", "text/css; charset=UTF-8")
			    		w.WriteHeader(http.StatusBadRequest)	
					} else if err.Error() == "401" {
						w.Header().Set("Content-Type", "text/css; charset=UTF-8")
			    		w.WriteHeader(http.StatusUnauthorized)	
					} else {
						hasexploded = err.Error()
					}		
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			    	w.WriteHeader(http.StatusOK)
				}   
			}
		} else {
			w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    	w.WriteHeader(http.StatusBadRequest)	
		}
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("Handlers.DeleteMessages() - HasExploded! - Error: %v", hasexploded))
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
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
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
				err := activityservice.CreateUserActivity(userId, activity)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
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
	
	password := r.Header.Get("key")
	
	if err := utilsservice.CheckUserPassword(userId, password); err == nil {
		if userId != "" {
		    response, err := activityservice.GetUserActivity(userId)
			
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
	} else if err.Error() == "401" || err.Error() == "not found" {
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusUnauthorized)
	} else {
		hasexploded = err.Error()
	}
	
	if hasexploded != "" {
		log.Printf(fmt.Sprintf("GetUserActivity - HasExploded! - Error: %v", hasexploded))
		w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	    w.WriteHeader(http.StatusInternalServerError)
	}
}

func RegisterInterest(w http.ResponseWriter, r *http.Request) {
	var interest types.InterestRequest
	var hasexploded string
	hasexploded = ""
	
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	
	if err != nil {
	    hasexploded = err.Error()
	} else if err := r.Body.Close(); err != nil {
	    hasexploded = err.Error()
	} else if err := json.Unmarshal(body, &interest); err != nil {
	    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	    w.WriteHeader(422) // unprocessable entity
	    if err := json.NewEncoder(w).Encode(err); err != nil {
	        hasexploded = err.Error()
	    }
	} else {
		err := utilsservice.RegisterInterest(interest)
		
		if err != nil {
			if err.Error() == "blank credentials" {
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

	if hasexploded != "" {
		log.Printf(fmt.Sprintf("Handlers.RegisterInterest() - HasExploded! - Error: %v", hasexploded))
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

func HelpPagePOSTUserLogin(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-user-login.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePOSTMatch(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-match.html")
    fmt.Fprint(w, string(body))
}

func HelpPageGETChats(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/GET-chats.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePUTChats(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/PUT-chats.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePOSTInterest(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-interest.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePUTExtraInfo(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/PUT-extrainfo.html")
    fmt.Fprint(w, string(body))
}

func HelpPagePOSTUserReport(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/POST-user-report.html")
    fmt.Fprint(w, string(body))
}

func HelpPageDELETEMessage(w http.ResponseWriter, r *http.Request) {
    body, _ := ioutil.ReadFile("helppage/DELETE-message.html")
    fmt.Fprint(w, string(body))
}

func ProduceKafkaMessage(w http.ResponseWriter, r *http.Request) {	
    kafkaservice.ProduceDummyMessage();
	w.Header().Set("Content-Type", "text/css; charset=UTF-8")
	w.WriteHeader(http.StatusOK)	
}

