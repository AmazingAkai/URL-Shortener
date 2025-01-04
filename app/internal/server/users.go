package server

import (
	"net/http"
	"strings"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database/queries"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"

	"github.com/go-chi/chi/v5"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userInput models.User

	if err := utils.ReadJSON(r.Body, &userInput); err != nil {
		utils.BadRequestError(w)
		return
	}
	if err := utils.ValidateStruct(userInput); err != nil {
		utils.ValidationError(w, err)
		return
	}

	user, err := queries.CreateUser(userInput)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\"") {
			utils.ErrorResponse(w, http.StatusConflict, "email already exists")
			return
		}
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}

func logInHandler(w http.ResponseWriter, r *http.Request) {
	var userInput models.User
	if err := utils.ReadJSON(r.Body, &userInput); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(userInput); err != nil {
		utils.ValidationError(w, err)
		return
	}

	user, err := queries.AuthenticateUser(userInput)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, exp, err := utils.GenerateJWT(user)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Map{
		"token": token,
		"exp":   exp,
	})
}

func (*Server) RegisterUserRoutes(r *chi.Mux) {
	r.Post("/register", createUserHandler)
	r.Post("/login", logInHandler)
}
