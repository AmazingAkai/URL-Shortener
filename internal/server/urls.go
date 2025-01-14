package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/internal/views"
	"github.com/AmazingAkai/URL-Shortener/internal/views/partials"

	"github.com/go-chi/chi/v5"
)

var (
	WEB_URL = os.Getenv("WEB_URL")
)

type UrlCreatePayload struct {
	ShortUrl string `schema:"short_url" validate:"required,min=5,max=30,alphanumeric,validShortUrl"`
	LongUrl  string `schema:"long_url" validate:"required,url,validLongUrl"`
}

func (s *Server) createShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var payload UrlCreatePayload
	if err := utils.ParseForm(r, &payload); err != nil {
		utils.ParseFormError(w, r, err)
		return
	}
	if err := utils.ValidateStruct(payload); err != nil {
		utils.ValidationError(w, r, err)
		return
	}

	url := &store.Url{
		ShortUrl: payload.ShortUrl,
		LongUrl:  payload.LongUrl,
	}
	session := s.getSession(r)
	if session != nil {
		url.UserID = &session.UserID
	}

	err := s.store.Urls.Create(r.Context(), url)
	if err != nil {
		switch err {
		case store.ErrConflict:
			utils.ErrorResponse(w, r, http.StatusConflict, []string{"Short URL already exists."})
		default:
			utils.ServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	partials.SuccessUrl(fmt.Sprintf("%s/%s", WEB_URL, url.ShortUrl)).Render(r.Context(), w)
}

func (s *Server) getShortUrlList(w http.ResponseWriter, r *http.Request) {
	session := s.getSession(r)
	if session == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, []string{"User not logged in."})
		return
	}

	urls, err := s.store.Urls.GetUrlList(r.Context(), session.UserID)
	if err != nil {
		utils.ServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	partials.UrlList(urls).Render(r.Context(), w)
}

func (s *Server) deleteShortUrl(w http.ResponseWriter, r *http.Request) {
	urlID, err := strconv.Atoi(chi.URLParam(r, "url_id"))
	if err != nil {
		utils.ErrorResponse(w, r, http.StatusBadRequest, []string{"Invalid URL ID."})
		return
	}

	session := s.getSession(r)
	if session == nil {
		utils.ErrorResponse(w, r, http.StatusUnauthorized, []string{"User not logged in."})
		return
	}

	if err = s.store.Urls.Delete(r.Context(), urlID, session.UserID); err != nil {
		switch err {
		case store.ErrNotFound:
			utils.ErrorResponse(w, r, http.StatusNotFound, []string{"URL not found."})
		default:
			utils.ServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) redirectShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	url, err := s.store.Urls.GetUrl(r.Context(), chi.URLParam(r, "short_url"))
	if err != nil {
		var (
			code    int
			message string
		)
		switch err {
		case store.ErrNotFound:
			code = 404
			message = "Oops! The page you're looking for doesn't exist or has been moved."
		default:
			code = 500
			message = "An internal server error occurred."
			log.Printf("Error getting long URL: %v", err)
		}

		w.WriteHeader(code)
		views.Layout(views.ErrorPage(strconv.Itoa(code), message)).Render(r.Context(), w)
		return
	}

	go s.store.Urls.IncrementVisits(url.ID)

	http.Redirect(w, r, url.LongUrl, http.StatusFound)
}
