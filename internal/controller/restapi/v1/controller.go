package v1

import (
	"github.com/andreyxaxa/Sales-Tracker/internal/usecase"
	"github.com/andreyxaxa/Sales-Tracker/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	c usecase.Category
	p usecase.Product
	s usecase.Sale
	a usecase.Analytics

	l logger.Interface

	v *validator.Validate
}
