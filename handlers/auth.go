package handlers

import (
	"auth-service/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var users = map[string]string{}    // In-memory user store (email -> hashed password)
var JwtKey = []byte("mysecretkey") // Exported secret key for JWT

// Struct for JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// SignUp handler to register a new user
func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	// Hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Store user in the map (in production, use a database)
	users[user.Email] = string(hashedPassword)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created successfully")
}

// SignIn handler to authenticate a user and generate a JWT token
func SignIn(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	storedPassword, exists := users[user.Email]
	if !exists || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tokenString)
}

// Revoked tokens list (in-memory store for simplicity)
var revokedTokens = map[string]bool{}

// RevokeToken handler to revoke a JWT token
func RevokeToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Add the token to the revoked list
	revokedTokens[tokenString] = true
	json.NewEncoder(w).Encode("Token revoked successfully")
}

// RefreshToken handler to refresh a valid JWT token before expiry
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if the token is about to expire
	if time.Until(time.Unix(claims.ExpiresAt, 0)) > 30*time.Minute {
		http.Error(w, "Token not close to expiry", http.StatusBadRequest)
		return
	}

	// Create a new token for the user with an updated expiration time
	expirationTime := time.Now().Add(1 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenString, err := newToken.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "Could not refresh token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newTokenString)
}
