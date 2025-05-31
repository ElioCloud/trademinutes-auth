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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.GetDB().Collection("MyClusterCol")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("JSON decode error:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Println("➡️ Incoming email:", user.Email)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"email": user.Email}).Err()
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Hashing failed", http.StatusInternalServerError)
		return
	}

	_, err = collection.InsertOne(ctx, user)
if err != nil {
    fmt.Println("❌ Insert error:", err) // 👈 this line
    http.Error(w, "Insert failed", http.StatusInternalServerError)
    return
}

	fmt.Println("User inserted:", user.Email)
	w.Write([]byte("User registered successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.GetDB().Collection("MyClusterCol")

	var input models.User
	var foundUser models.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&foundUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(input.Password, foundUser.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(foundUser.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.EmailKey).(string)
	collection := config.GetDB().Collection("MyClusterCol")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Never return password!
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}