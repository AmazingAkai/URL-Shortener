package server

import (
	"net/http"
	"time"

	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/internal/utils/constants"
)

type UserCreatePayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UserCreatePayload
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

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(constants.SESSION_KEY) != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "user already logged in")
		return
	}

	var payload UserCreatePayload
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

	sessionToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	csrfToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		utils.ServerError(w, err)
		return
	}
	expires := time.Now().Add(time.Hour * 24)

	session := &store.Session{
		UserID:    user.ID,
		Token:     sessionToken,
		CSRFToken: csrfToken,
		ExpiresAt: expires.Unix(),
	}
	s.store.Sessions.Set(session.Token, session)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expires,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  expires,
		HttpOnly: false,
	})

	utils.WriteJSON(w, http.StatusOK, utils.Map{
		"success": true,
	})
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	var session *store.Session
	ctx := r.Context()

	if ctx.Value(constants.SESSION_KEY) == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	session = ctx.Value(constants.SESSION_KEY).(*store.Session)
	s.store.Sessions.Delete(session.Token)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
	})

	utils.WriteJSON(w, http.StatusOK, utils.Map{
		"success": true,
	})
}
