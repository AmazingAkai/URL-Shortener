package server

import (
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/internal/middleware"
	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils"
)

type userCreatePayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload userCreatePayload
	if err := utils.ReadJSON(r.Body, &payload); err != nil {
		utils.BadRequestError(w)
		return
	}
	if err := utils.ValidateStruct(payload); err != nil {
		utils.ValidationError(w, err)
		return
	}

	user := &store.User{
		Email: payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		utils.ServerError(w, err)
		return
	}

	if err := s.store.Users.Create(r.Context(), user); err != nil {
		switch err {
		case store.ErrConflict:
			utils.ErrorResponse(w, http.StatusConflict, "user already exists")
		default:
			utils.ServerError(w, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, user)
}

func (s *Server) logInHandler(w http.ResponseWriter, r *http.Request) {
	var payload userCreatePayload
	if err := utils.ReadJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := utils.ValidateStruct(payload); err != nil {
		utils.ValidationError(w, err)
		return
	}

	user, err := s.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, exp, err := middleware.GenerateJWT(user)
	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Map{
		"token": token,
		"exp":   exp,
	})
}
