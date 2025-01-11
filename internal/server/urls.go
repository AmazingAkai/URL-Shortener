package server

import (
	"log"
	"net/http"
	"time"

	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/internal/utils/constants"

	"github.com/go-chi/chi/v5"
)

type UrlCreatePayload struct {
	LongUrl   string     `json:"long_url" validate:"required,url"`
	ExpiresAt *time.Time `json:"expires_at" validate:"omitempty,futureDate"`
}

func (s *Server) createShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var payload UrlCreatePayload
	if err := utils.ReadJSON(r.Body, &payload); err != nil {
		utils.ParseFormError(w, r, err)
		return
	}
	if err := utils.ValidateStruct(payload); err != nil {
		log.Printf("Validation error: %v", err)
		utils.ValidationError(w, r, err)
		return
	}

	url := &store.Url{
		LongUrl:   payload.LongUrl,
		ExpiresAt: payload.ExpiresAt,
	}
	session := r.Context().Value(constants.SESSION_KEY)
	if session != nil {
		url.UserID = &session.(*store.Session).UserID
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

	utils.WriteJSON(w, r, http.StatusCreated, url)
}

func (s *Server) redirectShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	url, err := s.store.Urls.GetLongUrl(r.Context(), chi.URLParam(r, "short_url"))
	if err != nil {
		switch err {
		case store.ErrNotFound:
			// utils.NotFoundError(w) // TODO: Fix this
			utils.ErrorResponse(w, r, http.StatusNotFound, []string{"Short URL not found."})
		default:
			utils.ServerError(w, r, err)
		}
		return
	}

	ipAddr := r.Header.Get("X-Forwarded-For")
	if ipAddr == "" {
		ipAddr = r.RemoteAddr
	}

	visit := store.UrlVisit{
		UrlID:     url.ID,
		IpAddr:    ipAddr,
		Referer:   r.Header.Get("Referer"),
		UserAgent: r.Header.Get("User-Agent"),
	}
	go s.store.Urls.CreateVisit(visit)

	http.Redirect(w, r, url.LongUrl, http.StatusFound)
}
