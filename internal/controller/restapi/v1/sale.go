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

// @Summary     Create sale
// @Description Creates a new sale with product id, quantity, unit price, payment method, and date of sale
// @Tags        sale
// @Accept		json
// @Produce		json
// @Param		request body request.CreateSaleRequest true "Product ID, Quantity, Unit price, Payment method, Sold at"
// @Success		201 {object} response.CreateSaleResponse
// @Failure		400 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/sale [post]
func (r *V1) createSale(ctx *fiber.Ctx) error {
	var body request.CreateSaleRequest

	err := ctx.BodyParser(&body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	err = r.v.Struct(body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	price, err := decimal.NewFromString(body.UnitPrice)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid price")
	}
	if price.LessThanOrEqual(decimal.Zero) {
		return errorResponse(ctx, fiber.StatusBadRequest, "price must be greater than 0")
	}

	if body.Quantity <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "quantity must be greater than 0")
	}

	if body.PaymentMethod != "cash" && body.PaymentMethod != "card" && body.PaymentMethod != "online" {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid payment method: must be one of: \"cash\" / \"card\" / \"online\"")
	}

	sl, err := r.s.Create(ctx.UserContext(), body.ProductID, body.Quantity, price, body.PaymentMethod, body.SoldAt)
	if err != nil {
		if errors.Is(err, errs.ErrProductNotFound) {
			return errorResponse(ctx, fiber.StatusBadRequest, "product nto found")
		}
		r.l.Error(err, "restapi - v1 - createSale")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := response.CreateSaleResponse{
		ID:            sl.ID,
		ProductID:     sl.ProductID,
		Quantity:      sl.Quantity,
		UnitPrice:     sl.UnitPrice.String(),
		PaymentMethod: sl.PaymentMethod,
		SoldAt:        sl.SoldAt,
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// @Summary     Get sale
// @Description Returns sale by ID
// @Tags        sale
// @Produce		json
// @Param		id path string true "Sale ID"
// @Success		200 {object} response.GetSaleResponse
// @Failure		400 {object} response.Error
// @Failure		404 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/sale/{id} [get]
func (r *V1) getSale(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	if id <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id: id cant be 0 or negative")
	}

	sl, err := r.s.Get(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, errs.ErrSaleNotFound) {
			return errorResponse(ctx, fiber.StatusNotFound, "sale not found")
		}
		r.l.Error(err, "restapi - v1 - getSale")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := response.GetSaleResponse{
		ID:            sl.ID,
		ProductID:     sl.ProductID,
		Quantity:      sl.Quantity,
		UnitPrice:     sl.UnitPrice.String(),
		PaymentMethod: sl.PaymentMethod,
		SoldAt:        sl.SoldAt,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary     Delete sale
// @Description Deletes sale by ID
// @Tags        sale
// @Param		id path string true "Sale ID"
// @Success		204 "Deleted"
// @Failure		400 {object} response.Error
// @Failure		404 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/sale/{id} [delete]
func (r *V1) deleteSale(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	if id <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id: id cant be 0 or negative")
	}

	err = r.s.Delete(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, errs.ErrSaleNotFound) {
			return errorResponse(ctx, fiber.StatusNotFound, "sale not found")
		}
		r.l.Error(err, "restapi - v1 - deleteSale")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
