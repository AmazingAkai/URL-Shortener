package server

import (
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/internal/views"
)

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	var (
		urls []*store.Url
		err  error
	)
	if session := s.getSession(r); session != nil {
		urls, err = s.store.Urls.GetUrlList(r.Context(), session.UserID)
		if err != nil {
			utils.ServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	views.Layout(views.Home(urls)).Render(r.Context(), w)
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	views.Layout(views.ErrorPage("404", "Oops! The page you're looking for doesn't exist or has been moved.")).Render(r.Context(), w)
}

func (s *Server) registerPageHandler(w http.ResponseWriter, r *http.Request) {
	if session := s.getSession(r); session != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	views.Layout(views.Register()).Render(r.Context(), w)
}

func (s *Server) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if session := s.getSession(r); session != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	views.Layout(views.Login()).Render(r.Context(), w)
}
