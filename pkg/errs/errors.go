package errs

import "errors"

var (
	// Category
	ErrCategoryAlreadyExists = errors.New("category with this name already exists")
	ErrCategoryHasProducts   = errors.New("there are still products in this category")
	ErrCategoryNotFound      = errors.New("category not found")

	// Product
	ErrProductNotFound = errors.New("product not found")

	// Sale
	ErrSaleNotFound = errors.New("sale not found")
)
