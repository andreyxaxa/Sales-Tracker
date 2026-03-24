package restapi

import (
	"github.com/andreyxaxa/Sales-Tracker/config"
	_ "github.com/andreyxaxa/Sales-Tracker/docs" // Swagger docs.
	v1 "github.com/andreyxaxa/Sales-Tracker/internal/controller/restapi/v1"
	"github.com/andreyxaxa/Sales-Tracker/internal/usecase"
	"github.com/andreyxaxa/Sales-Tracker/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title    Sales Tracker
// @version  1.0.0
// @host     localhost:8080
// @BasePath /v1
func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	c usecase.Category,
	p usecase.Product,
	s usecase.Sale,
	a usecase.Analytics,
	l logger.Interface,
) {
	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewRoutes(apiV1Group, c, p, s, a, l)
	}
}
