package routes

import (
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database/queries"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
	"github.com/gorilla/mux"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := utils.ReadJSON(r.Body, &user); err != nil {
		utils.BadRequestError(w)
		return
	}

	if err := utils.ValidateStruct(user); err != nil {
		utils.ValidationError(w, err)
		return
	}

	if err := queries.CreateUser(&user); err != nil {
		if err.Error() == "email already exists" {
			utils.ErrorResponse(w, http.StatusConflict, err.Error())
			return
		}
		utils.ServerError(w, err)
		return
	}

	user.Password = "" // Don't return the password
	utils.WriteJSON(w, http.StatusCreated, user)
}

func logInHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := utils.ReadJSON(r.Body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(user); err != nil {
		utils.ValidationError(w, err)
		return
	}

	if err := queries.AuthenticateUser(&user); err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, exp, err := utils.GenerateJWT(&user)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Map{
		"token": token,
		"exp":   exp,
	})
}

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/register", createUserHandler).Methods("POST")
	router.HandleFunc("/login", logInHandler).Methods("POST")
}
