package routes

import (
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database/queries"
	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils/constants"

	"github.com/gorilla/mux"
)

func redirectShortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	longURL, err := queries.GetLongURL(vars["short_url"])
	if err != nil || longURL == "" {
		utils.NotFoundError(w)
		return
	}

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

func RegisterURLRoutes(router *mux.Router) {
	router.HandleFunc("/{short_url}", redirectShortURLHandler).Methods("GET")
	router.HandleFunc("/urls", createShortURLHandler).Methods("POST")
}
