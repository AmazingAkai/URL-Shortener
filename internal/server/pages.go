package server

import (
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/internal/views"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	session := s.getSession(r)
	views.Home(session).Render(r.Context(), w)
}

func (s *Server) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	session := s.getSession(r)
	if session != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	views.Login().Render(r.Context(), w)
}
