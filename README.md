
# TradeMinutes Auth API

This is a simple user authentication API built with **Go** and **MongoDB**, using **JWT** for session management. It includes registration, login, and protected routes.

---

## 🛠️ Features

- Register new users
- Secure login with JWT
- Password hashing using bcrypt
- JWT-based route protection
- MongoDB for user storage

---

## ⚙️ Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/ElioCloud/trademinutes-auth
cd trademinutes-auth

```

Initialize Go module

```bash
go mod init trademinutes-auth
go mod tidy
```

Create a .env file in the root directory with the following:

```bash
# App Config
PORT=8080

# MongoDB Atlas
MONGO_URI=
DB_NAME=

# JWT Secrets
JWT_SECRET=
JWT_RESET_SECRET=

# Mailtrap SMTP for Password Reset Emails
SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASS=
EMAIL_FROM=

# Frontend URL (for reset password links)
FRONTEND_URL=

```

Run the server

```bash
go run main.go
```

## 🌐 Base URLs

| Environment | URL |
|-------------|-----|
| Local       | `http://localhost:8080` |
| Production  | `https://trademinutes-auth.onrender.com` |

---

## 📮 API Endpoints

### ➕ `POST /api/auth/register`

Registers a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "secret123",
  "name": "Test User"
}
```

### ➕ `POST /api/auth/login`

Login for the user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "secret123"
}

```

### ➕ `POST /api/auth/forgot-password`

Forgot the password.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

### ➕ `POST /api/auth/reset-password`

Forgot the password.

**Request Body:**
```json
{
  "token": "<token>",
  "newPassword": "new123456"
}


