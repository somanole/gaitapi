// GaitAPI
package main

import (
    "log"
    "net/http"
)

func main() {

    router := NewRouter()

    log.Fatal(http.ListenAndServe("127.0.0.1:3000", router))
}


