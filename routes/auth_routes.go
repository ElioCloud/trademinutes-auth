package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"trademinutes-auth/controllers"
)

func AuthRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/api/auth").Subrouter()

	authRouter.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")
	authRouter.HandleFunc("/login", controllers.LoginHandler).Methods("POST")

}
