package routes

import (
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database/queries"
	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils/constants"

	"github.com/go-chi/chi/v5"
)

func redirectShortURLHandler(w http.ResponseWriter, r *http.Request) {
	id, longURL, err := queries.GetLongURL(chi.URLParam(r, "short_url"))
	if err != nil || longURL == "" {
		utils.NotFoundError(w)
		return
	}

	ipAddr := r.Header.Get("X-Forwarded-For")
	if ipAddr == "" {
		ipAddr = r.RemoteAddr
	}
	go queries.CreateVisit(id, ipAddr, r.Header.Get("Referer"), r.Header.Get("User-Agent"))

	http.Redirect(w, r, longURL, http.StatusFound)
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var urlInput models.URL

	if err := utils.ReadJSON(r.Body, &urlInput); err != nil {
		utils.BadRequestError(w)
		return
	}
	if err := utils.ValidateStruct(urlInput); err != nil {
		log.Errorf("Validation error: %v", err)
		utils.ValidationError(w, err)
		return
	}

	user := r.Context().Value(constants.USER_KEY)
	url, err := queries.CreateShortURL(urlInput, user)

	if err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, url)
}

func RegisterURLRoutes(r *chi.Mux) {
	r.Get("/{short_url}", redirectShortURLHandler)
	r.Post("/urls", createShortURLHandler)
}
