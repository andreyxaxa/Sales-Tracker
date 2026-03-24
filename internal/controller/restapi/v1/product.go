package v1

import (
	"errors"
	"strconv"

	"github.com/andreyxaxa/Sales-Tracker/internal/controller/restapi/v1/request"
	"github.com/andreyxaxa/Sales-Tracker/internal/controller/restapi/v1/response"
	"github.com/andreyxaxa/Sales-Tracker/pkg/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

// @Summary     Create product
// @Description Creates a new product with category id, name and base price
// @Tags        product
// @Accept		json
// @Produce		json
// @Param		request body request.CreateProductRequest true "Category ID, Name, Base price"
// @Success		201 {object} response.CreateProductResponse
// @Failure		400 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/product [post]
func (r *V1) createProduct(ctx *fiber.Ctx) error {
	var body request.CreateProductRequest

	err := ctx.BodyParser(&body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	err = r.v.Struct(body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	basePrice, err := decimal.NewFromString(body.BasePrice)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid base price")
	}

	if basePrice.LessThanOrEqual(decimal.Zero) {
		return errorResponse(ctx, fiber.StatusBadRequest, "base price must be positive and greater than 0")
	}

	prd, err := r.p.Create(ctx.UserContext(), body.CategoryID, body.Name, basePrice)
	if err != nil {
		if errors.Is(err, errs.ErrCategoryNotFound) {
			return errorResponse(ctx, fiber.StatusBadRequest, "category not found")
		}
		r.l.Error(err, "restapi - v1 - createProduct")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := response.CreateProductResponse{
		ID:         prd.ID,
		CategoryID: prd.CategoryID,
		Name:       prd.Name,
		BasePrice:  prd.BasePrice.String(),
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// @Summary     Get product
// @Description Returns product by ID
// @Tags        product
// @Produce		json
// @Param		id path string true "Produt ID"
// @Success		200 {object} response.GetProductResponse
// @Failure		400 {object} response.Error
// @Failure		404 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/product/{id} [get]
func (r *V1) getProduct(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	if id <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id: id cant be 0 or negative")
	}

	prd, err := r.p.Get(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			return errorResponse(ctx, fiber.StatusNotFound, "product not found")
		}
		r.l.Error(err, "restapi - v1 - getProduct")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := response.GetProductResponse{
		ID:         prd.ID,
		CategoryID: prd.CategoryID,
		Name:       prd.Name,
		BasePrice:  prd.BasePrice.String(),
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary     Delete product
// @Description Deletes product by ID
// @Tags        product
// @Param		id path string true "Produt ID"
// @Success		204 "Deleted"
// @Failure		400 {object} response.Error
// @Failure		404 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/product/{id} [delete]
func (r *V1) deleteProduct(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	if id <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id: id cant be 0 or negative")
	}

	err = r.p.Delete(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			return errorResponse(ctx, fiber.StatusNotFound, "product not found")
		}
		r.l.Error(err, "restapi - v1 - deleteProduct")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary     Update price of product
// @Description Updates product's price
// @Tags        product
// @Accept		json
// @Produce		json
// @Param		id path string true "Product ID"
// @Param		request body request.UpdateProductPriceRequest true "New price"
// @Success		200 {object} response.UpdateProductBasePriceResponse
// @Failure		400 {object} response.Error
// @Failure		404 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/product/{id} [patch]
func (r *V1) updateProductBasePrice(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	if id <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id: id cant be 0 or negative")
	}

	var body request.UpdateProductPriceRequest

	err = ctx.BodyParser(&body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	err = r.v.Struct(body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	newPrice, err := decimal.NewFromString(body.NewPrice)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid price")
	}

	if newPrice.LessThanOrEqual(decimal.Zero) {
		return errorResponse(ctx, fiber.StatusBadRequest, "price must be positive and greater than 0")
	}

	updatedProduct, err := r.p.UpdateBasePrice(ctx.UserContext(), id, newPrice)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			return errorResponse(ctx, fiber.StatusNotFound, "product not found")
		}
		r.l.Error(err, "restapi - v1 - updateProductBasePrice")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := response.UpdateProductBasePriceResponse{
		ID:         updatedProduct.ID,
		CategoryID: updatedProduct.CategoryID,
		Name:       updatedProduct.Name,
		BasePrice:  updatedProduct.BasePrice.String(),
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
