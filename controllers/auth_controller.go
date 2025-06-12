package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"trademinutes-auth/config"
	"trademinutes-auth/models"
	"trademinutes-auth/utils"
	"trademinutes-auth/middleware"
	"go.mongodb.org/mongo-driver/bson"
)

// Helper to send consistent JSON error responses
func writeJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.GetDB().Collection("MyClusterCol")

	var input models.User
	var foundUser models.User

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		writeJSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Lookup user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&foundUser)
	if err != nil {
		writeJSONError(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(input.Password, foundUser.Password) {
		writeJSONError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := utils.GenerateJWT(foundUser.Email)
	if err != nil {
		writeJSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// ✅ Return token in JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
