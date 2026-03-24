package request

type CreateProductRequest struct {
	CategoryID int64  `json:"category_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
	BasePrice  string `json:"base_price" validate:"required"`
}

type UpdateProductPriceRequest struct {
	NewPrice string `json:"new_price" validate:"required"`
}
