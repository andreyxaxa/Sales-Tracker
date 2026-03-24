package v1

import (
	"errors"
	"strconv"

	"github.com/andreyxaxa/Sales-Tracker/internal/controller/restapi/v1/request"
	"github.com/andreyxaxa/Sales-Tracker/internal/controller/restapi/v1/response"
	"github.com/andreyxaxa/Sales-Tracker/pkg/errs"
	"github.com/gofiber/fiber/v2"
)

// @Summary     Create category
// @Description Creates a new category with name and description(optionally)
// @Tags        category
// @Accept		json
// @Produce		json
// @Param		request body request.CreateCategoryRequest true "Name, description(opt.)"
// @Success		201 {object} response.CreateCategoryResponse
// @Failure		400 {object} response.Error
// @Failure		409 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/category [post]
func (r *V1) createCategory(ctx *fiber.Ctx) error {
	var body request.CreateCategoryRequest

	err := ctx.BodyParser(&body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	err = r.v.Struct(body)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid request body")
	}

	ctg, err := r.c.Create(ctx.UserContext(), body.Name, body.Description)
	if err != nil {
		if errors.Is(err, errs.ErrCategoryAlreadyExists) {
			return errorResponse(ctx, fiber.StatusConflict, "category already exists")
		}
		r.l.Error(err, "restapi - v1 - createCategory")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := response.CreateCategoryResponse{
		ID:          ctg.ID,
		Name:        ctg.Name,
		Description: ctg.Description,
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

// @Summary     Get category
// @Description Returns category by ID
// @Tags        category
// @Produce		json
// @Param		id path string true "Category ID"
// @Success		200 {object} response.GetCategoryResponse
// @Failure		400 {object} response.Error
// @Failure		404 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/category/{id} [get]
func (r *V1) getCategoryByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	if id <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id: id cant be 0 or negative")
	}

	ctg, err := r.c.GetByID(ctx.UserContext(), id)
	if err != nil {
		if errors.Is(err, errs.ErrCategoryNotFound) {
			return errorResponse(ctx, fiber.StatusNotFound, "category not found")
		}
		r.l.Error(err, "restapi - v1 - getCategoryByID")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := response.GetCategoryResponse{
		ID:          ctg.ID,
		Name:        ctg.Name,
		Description: ctg.Description,
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary     Delete
// @Description Deletes category by ID
// @Tags        category
// @Param		id path string true "Category ID"
// @Success		204 "Deleted"
// @Failure		400 {object} response.Error
// @Failure		404 {object} response.Error
// @Failure		409 {object} response.Error
// @Failure		500 {object} response.Error
// @Router 		/v1/category/{id} [delete]
func (r *V1) deleteCategoryByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id")
	}

	if id <= 0 {
		return errorResponse(ctx, fiber.StatusBadRequest, "invalid id: id cant be 0 or negative")
	}

	err = r.c.DeleteByID(ctx.UserContext(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrCategoryHasProducts):
			return errorResponse(ctx, fiber.StatusConflict, "category has products")
		case errors.Is(err, errs.ErrCategoryNotFound):
			return errorResponse(ctx, fiber.StatusNotFound, "category not found")
		default:
			r.l.Error(err, "restapi - v1 - deleteCategoryByID")
			return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
		}
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
