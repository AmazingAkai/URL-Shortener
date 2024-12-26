package routes

import (
	"github.com/AmazingAkai/URL-Shortener/app/internal/database/queries"
	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/server"
	"github.com/AmazingAkai/URL-Shortener/app/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func RedirectShortURLHandler(c *fiber.Ctx) error {
	shortURL := c.Params("short_url")

	longURL, err := queries.GetLongURL(shortURL)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Short URL not found")
	}

	return c.Redirect(longURL)
}

func CreateShortURLHandler(c *fiber.Ctx) error {
	var url models.URL
	if err := c.BodyParser(&url); err != nil {
		return err
	}

	url.ShortURL = utils.GenerateShortURL()
	// TODO: Check if this short url already exists in the database

	if err := validate.Struct(url); err != nil {
		return err
	}

	if err := queries.CreateShortURL(url); err != nil {
		return err
	}

	return c.JSON(url)
}

func RegisterURLRoutes(s *server.FiberServer) {
	s.App.Get("/:short_url", RedirectShortURLHandler)

	api := s.App.Group("/urls")
	api.Post("/", CreateShortURLHandler)
}
