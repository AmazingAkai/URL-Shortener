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

func (s *Server) redirectShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	url, err := s.store.Urls.GetLongUrl(r.Context(), chi.URLParam(r, "short_url"))
	if err != nil {
		switch err {
		case store.ErrNotFound:
			utils.NotFoundError(w)
		default:
			utils.ServerError(w, err)
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
	go s.store.Urls.CreateVisit(r.Context(), visit)

	http.Redirect(w, r, url.LongUrl, http.StatusFound)
}

func (s *Server) createShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var payload UrlCreatePayload

	if err := utils.ReadJSON(r.Body, &payload); err != nil {
		utils.BadRequestError(w)
		return
	}
	if err := utils.ValidateStruct(payload); err != nil {
		log.Printf("Validation error: %v", err)
		utils.ValidationError(w, err)
		return
	}

	url := &store.Url{
		LongUrl:   payload.LongUrl,
		ExpiresAt: payload.ExpiresAt,
	}
	user := r.Context().Value(constants.USER_KEY)
	if user != nil {
		url.UserID = &user.(*store.User).ID
	}

	err := s.store.Urls.Create(r.Context(), url)
	if err != nil {
		switch err {
		case store.ErrConflict:
			utils.ErrorResponse(w, http.StatusConflict, "short url already exists")
		default:
			utils.ServerError(w, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, url)
}
