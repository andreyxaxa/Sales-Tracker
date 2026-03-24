package response

type CategoryResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateCategoryResponse CategoryResponse
type GetCategoryResponse CategoryResponse
