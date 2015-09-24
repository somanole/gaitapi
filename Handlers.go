// Handlers
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
	"io"
	"io/ioutil"
	"github.com/somanole/gaitapi/acceleration"
	"github.com/somanole/gaitapi/repo"
	"github.com/somanole/gaitapi/services"
)

var repository repo.Repo

func init() {
	repository = repo.New()
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome Sorin!")
}

func AccelerationIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	
	accelerations := repository.GetAllAccelerations()
    if err := json.NewEncoder(w).Encode(accelerations); err != nil {
        panic(err)
    }
}

func AccelerationsCount(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	
	count := repository.GetAccelerationsCount()
    if err := json.NewEncoder(w).Encode(count); err != nil {
        panic(err)
    }
}

func AccelerationCreate(w http.ResponseWriter, r *http.Request) {
	var acceleration acceleration.Acceleration
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

    a := repository.CreateAcceleration(acceleration)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(a); err != nil {
        panic(err)
    }
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
