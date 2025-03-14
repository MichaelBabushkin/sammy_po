package main

import (
    "log"
    "net/http"
)

func main() {
    // Serve files in the "frontend" directory at the root path.
    fs := http.FileServer(http.Dir("frontend"))
    http.Handle("/", fs)

    log.Println("Starting server on :8080...")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
