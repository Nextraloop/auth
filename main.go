package main

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/signup", handlers.SignUp).Methods("POST")
	r.HandleFunc("/signin", handlers.SignIn).Methods("POST")
	r.HandleFunc("/revoke", handlers.RevokeToken).Methods("POST")
	r.HandleFunc("/refresh", handlers.RefreshToken).Methods("POST")

	// Protected route example
	r.Handle("/protected", middleware.AuthMiddleware(http.HandlerFunc(ProtectedHandler)))

	http.ListenAndServe(":8080", r)
}

// ProtectedHandler for testing token authorization
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Access granted to protected resource")
}
