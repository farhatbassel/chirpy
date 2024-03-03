package main

import (
	"fmt"
	"net/http"
)

func main() {
    mux := http.NewServeMux()

    mux.Handle("/", http.FileServer(http.Dir(".")))

    corsMux := middlewareCors(mux)

    fmt.Println("Starting server on port 8080...")
    server := http.Server{
        Addr: "localhost:8080",
        Handler: corsMux,
    }
    
    server.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "*")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}
