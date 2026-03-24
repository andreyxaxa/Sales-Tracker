package v1

import (
	"github.com/andreyxaxa/Sales-Tracker/internal/usecase"
	"github.com/andreyxaxa/Sales-Tracker/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(
	apiV1Group fiber.Router,
	c usecase.Category,
	p usecase.Product,
	s usecase.Sale,
	a usecase.Analytics,
	l logger.Interface,
) {
	r := &V1{
		c: c,
		p: p,
		s: s,
		a: a,
		l: l,
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	categoryGroup := apiV1Group.Group("/category")
	productGroup := apiV1Group.Group("/product")
	saleGroup := apiV1Group.Group("/sale")
	analyticsGroup := apiV1Group.Group("/analytics")

	{
		// UI
		apiV1Group.Get("/", r.showUI)

		// API
		categoryGroup.Post("/", r.createCategory)
		categoryGroup.Get("/:id", r.getCategoryByID)
		categoryGroup.Delete("/:id", r.deleteCategoryByID)

		productGroup.Post("/", r.createProduct)
		productGroup.Get("/:id", r.getProduct)
		productGroup.Patch("/:id", r.updateProductBasePrice)
		productGroup.Delete("/:id", r.deleteProduct)

		saleGroup.Post("/", r.createSale)
		saleGroup.Get("/:id", r.getSale)
		saleGroup.Delete(":id", r.deleteSale)

		analyticsGroup.Get("/revenue-sum", r.getRevenueSum)
		analyticsGroup.Get("/revenue-avg", r.getRevenueAvg)
		analyticsGroup.Get("/sales-count", r.getSalesCount)
		analyticsGroup.Get("/median", r.getMedian)
		analyticsGroup.Get("/p90", r.getPercentile90)
	}
}
