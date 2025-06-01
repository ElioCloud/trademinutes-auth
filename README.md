
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
MONGO_URI=mongodb+srv://mathewsteven1996:steven123@mycluster.wxbgwit.mongodb.net/?retryWrites=true&w=majority&appName=MyCluster
DB_NAME=MyClusterDB

# JWT Secrets
JWT_SECRET=your_super_secret
JWT_RESET_SECRET=your_reset_secret

# Mailtrap SMTP for Password Reset Emails
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=dd4b13d9c51309
SMTP_PASS=84a203810ee0fc
EMAIL_FROM=no-reply@trademinutes.com

# Frontend URL (for reset password links)
FRONTEND_URL=http://localhost:3000

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


