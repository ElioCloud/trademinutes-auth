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
