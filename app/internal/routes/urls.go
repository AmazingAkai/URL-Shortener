package routes

import (
	"fmt"
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/app/internal/database/queries"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
	"github.com/gorilla/mux"
)

func redirectShortURLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	longURL, err := queries.GetLongURL(vars["short_url"])
	if err != nil {
		utils.NotFoundError(w)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var url models.URL
	if err := utils.ReadJSON(r.Body, &url); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	attempts := 0
	for {
		if attempts >= 10 {
			utils.ServerError(w, fmt.Errorf("failed to generate unique short URL after %d attempts", attempts))
			return
		}

		url.ShortURL = utils.GenerateShortURL()
		exists, err := queries.ShortURLExists(url.ShortURL)
		if err != nil {
			utils.ServerError(w, err)
			return
		}

		if !exists {
			break
		}
		attempts++
	}

	if err := utils.ValidateStruct(url); err != nil {
		utils.ValidationError(w, err)
		return
	}
	if err := queries.CreateShortURL(url); err != nil {
		utils.ServerError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, url)
}

func RegisterURLRoutes(router *mux.Router) {
	router.HandleFunc("/{short_url}/", redirectShortURLHandler).Methods("GET")
	router.HandleFunc("/urls/", createShortURLHandler).Methods("POST")
}
