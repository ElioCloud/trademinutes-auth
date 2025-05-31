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

	"go.mongodb.org/mongo-driver/bson"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	collection := config.GetDB().Collection("MyClusterCol") // ✅ moved inside handler

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("❌ JSON decode error:", err)
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
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ User inserted:", user.Email)
	w.Write([]byte("User registered successfully"))
}
