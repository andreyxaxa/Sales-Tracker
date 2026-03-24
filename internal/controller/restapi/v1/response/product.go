package response

type ProductResponse struct {
	ID         int64  `json:"id"`
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
	BasePrice  string `json:"base_price"`
}

type CreateProductResponse ProductResponse
type GetProductResponse ProductResponse
type UpdateProductBasePriceResponse ProductResponse
