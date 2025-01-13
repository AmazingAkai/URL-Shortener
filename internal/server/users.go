package server

import (
	"net/http"
	"time"

	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/internal/utils/constants"
)

type UserCreatePayload struct {
	Email    string `schema:"email" validate:"required,email"`
	Password string `schema:"password" validate:"required,min=8,max=32"`
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload UserCreatePayload
	if err := utils.ParseForm(r, &payload); err != nil {
		utils.ParseFormError(w, r, err)
		return
	}
	if err := utils.ValidateStruct(payload); err != nil {
		utils.ValidationError(w, r, err)
		return
	}

	user := &store.User{
		Email: payload.Email,
	}
	if err := user.Password.Set(payload.Password); err != nil {
		utils.ServerError(w, r, err)
		return
	}

	if err := s.store.Users.Create(r.Context(), user); err != nil {
		switch err {
		case store.ErrConflict:
			utils.ErrorResponse(w, r, http.StatusConflict, []string{"User with that email already exists."})
		default:
			utils.ServerError(w, r, err)
		}
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(constants.SESSION_KEY) != nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, []string{"User already logged in."})
		return
	}

	var payload UserCreatePayload
	if err := utils.ParseForm(r, &payload); err != nil {
		utils.ParseFormError(w, r, err)
		return
	}

	if err := utils.ValidateStruct(payload); err != nil {
		utils.ValidationError(w, r, err)
		return
	}

	user, err := s.store.Users.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.ErrorResponse(w, r, http.StatusUnauthorized, []string{"Invalid email or password."})
		default:
			utils.ServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, []string{"Invalid email or password."})
		return
	}

	sessionToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		utils.ServerError(w, r, err)
		return
	}

	expires := time.Now().Add(time.Hour * 24)
	session := &store.Session{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: expires.Unix(),
	}
	s.store.Sessions.Set(session.Token, session)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expires,
		HttpOnly: true,
	})

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session := s.getSession(r)
	if session == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	s.store.Sessions.Delete(session.Token)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) getSession(r *http.Request) *store.Session {
	if session, ok := r.Context().Value(constants.SESSION_KEY).(*store.Session); ok {
		return session
	}
	return nil
}
