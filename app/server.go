package main

import (
    "fmt"
    "net/http"
    "os"
)

func helloWorld(w http.ResponseWriter, r *http.Request){
    name, err := os.Hostname()

    if err != nil {
        panic(err)
    }
    
    fmt.Fprintf(w, "Hello World!!!!")
    fmt.Fprintf(w, name)
}

func main() {
    http.HandleFunc("/", helloWorld)
    http.ListenAndServe(":8080", nil)
}
