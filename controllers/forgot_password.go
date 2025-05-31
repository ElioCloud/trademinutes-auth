package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"trademinutes-auth/config"
	"trademinutes-auth/models"
)

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if user exists
	collection := config.GetDB().Collection("MyClusterCol")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate reset token
	resetToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	})
	secret := []byte(os.Getenv("JWT_RESET_SECRET"))
	tokenString, _ := resetToken.SignedString(secret)

	// Send email
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), tokenString)
	sendEmail(req.Email, "Password Reset", fmt.Sprintf("Click to reset password: %s", resetURL))

	w.Write([]byte("Password reset email sent"))
}

func sendEmail(to, subject, body string) {
	from := os.Getenv("EMAIL_FROM")
	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), os.Getenv("SMTP_HOST"))

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	_ = smtp.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, from, []string{to}, msg)
}


func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" || req.NewPassword == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Parse the token
	secret := []byte(os.Getenv("JWT_RESET_SECRET"))
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["email"] == nil {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	email := claims["email"].(string)
	collection := config.GetDB().Collection("MyClusterCol")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}

	_, err = collection.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": bson.M{"password": hashedPassword}},
	)

	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Password updated successfully"))
}
