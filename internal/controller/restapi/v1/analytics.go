package v1

import (
	"strconv"
	"time"

	"github.com/andreyxaxa/Sales-Tracker/internal/controller/restapi/v1/response"
	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/gofiber/fiber/v2"
)

// @Summary     Get revenue sum
// @Description Returns the revenue sum for period
// @Tags        analytics
// @Produce 	json
// @Param       from           query string true "Start of range (from=2026-01-10T12:15:31+03:00)"
// @Param 		to             query string true "End of range   (to=2026-03-22T19:07:13+03:00)"
// @Param 		product-id     query string false "If you need a specific product"
// @Param 		category-id    query string false "If you need a specific category"
// @Param 		payment-method query string false "If you need a specific payment method"
// @Param 		group-by       query string false "Grouping by different fields (group-by=day/month/week/category/product)"
// @Param 		sort-by        query string false "Sorting by field (sort-by=group/value)"
// @Param 		sort-order     query string false "Sorting order (sort-order=asc/desc)"
// @Success 	200 {object} []response.AnalyticsResponse
// @Failure 	400 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/analytics/revenue-sum [get]
func (r *V1) getRevenueSum(ctx *fiber.Ctx) error {
	params, s := parseAnalyticsParams(ctx)
	if s != "" {
		return errorResponse(ctx, fiber.StatusBadRequest, s)
	}

	results, err := r.a.GetRevenue(ctx.UserContext(), params)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getRevenueSum")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := make([]response.AnalyticsResponse, len(results))
	for i, r := range results {
		resp[i] = response.AnalyticsResponse{
			GroupLabel: r.GroupLabel,
			Value:      r.Value,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary     Get revenue average
// @Description Returns the average revenue for period
// @Tags        analytics
// @Produce 	json
// @Param       from           query string true "Start of range (from=2026-01-10T12:15:31+03:00)"
// @Param 		to             query string true "End of range   (to=2026-03-22T19:07:13+03:00)"
// @Param 		product-id     query string false "If you need a specific product"
// @Param 		category-id    query string false "If you need a specific category"
// @Param 		payment-method query string false "If you need a specific payment method"
// @Param 		group-by       query string false "Grouping by different fields (group-by=day/month/week/category/product)"
// @Param 		sort-by        query string false "Sorting by field (sort-by=group/value)"
// @Param 		sort-order     query string false "Sorting order (sort-order=asc/desc)"
// @Success 	200 {object} []response.AnalyticsResponse
// @Failure 	400 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/analytics/revenue-avg [get]
func (r *V1) getRevenueAvg(ctx *fiber.Ctx) error {
	params, s := parseAnalyticsParams(ctx)
	if s != "" {
		return errorResponse(ctx, fiber.StatusBadRequest, s)
	}

	results, err := r.a.GetAvgRevenue(ctx.UserContext(), params)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getRevenueAvg")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := make([]response.AnalyticsResponse, len(results))
	for i, r := range results {
		resp[i] = response.AnalyticsResponse{
			GroupLabel: r.GroupLabel,
			Value:      r.Value,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary     Get sales count
// @Description Returns the number of sales for period
// @Tags        analytics
// @Produce 	json
// @Param       from           query string true "Start of range (from=2026-01-10T12:15:31+03:00)"
// @Param 		to             query string true "End of range   (to=2026-03-22T19:07:13+03:00)"
// @Param 		product-id     query string false "If you need a specific product"
// @Param 		category-id    query string false "If you need a specific category"
// @Param 		payment-method query string false "If you need a specific payment method"
// @Param 		group-by       query string false "Grouping by different fields (group-by=day/month/week/category/product)"
// @Param 		sort-by        query string false "Sorting by field (sort-by=group/value)"
// @Param 		sort-order     query string false "Sorting order (sort-order=asc/desc)"
// @Success 	200 {object} []response.AnalyticsResponse
// @Failure 	400 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/analytics/sales-count [get]
func (r *V1) getSalesCount(ctx *fiber.Ctx) error {
	params, s := parseAnalyticsParams(ctx)
	if s != "" {
		return errorResponse(ctx, fiber.StatusBadRequest, s)
	}

	results, err := r.a.GetSalesCount(ctx.UserContext(), params)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getSalesCount")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := make([]response.AnalyticsResponse, len(results))
	for i, r := range results {
		resp[i] = response.AnalyticsResponse{
			GroupLabel: r.GroupLabel,
			Value:      r.Value,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary     Get median
// @Description Returns the median(p50) for period
// @Tags        analytics
// @Produce 	json
// @Param       from           query string true "Start of range (from=2026-01-10T12:15:31+03:00)"
// @Param 		to             query string true "End of range   (to=2026-03-22T19:07:13+03:00)"
// @Param 		product-id     query string false "If you need a specific product"
// @Param 		category-id    query string false "If you need a specific category"
// @Param 		payment-method query string false "If you need a specific payment method"
// @Param 		group-by       query string false "Grouping by different fields (group-by=day/month/week/category/product)"
// @Param 		sort-by        query string false "Sorting by field (sort-by=group/value)"
// @Param 		sort-order     query string false "Sorting order (sort-order=asc/desc)"
// @Success 	200 {object} []response.AnalyticsResponse
// @Failure 	400 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/analytics/median [get]
func (r *V1) getMedian(ctx *fiber.Ctx) error {
	params, s := parseAnalyticsParams(ctx)
	if s != "" {
		return errorResponse(ctx, fiber.StatusBadRequest, s)
	}

	results, err := r.a.GetMedian(ctx.UserContext(), params)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getMedian")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := make([]response.AnalyticsResponse, len(results))
	for i, r := range results {
		resp[i] = response.AnalyticsResponse{
			GroupLabel: r.GroupLabel,
			Value:      r.Value,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// @Summary     Get percentile 90
// @Description Returns the p90 for period
// @Tags        analytics
// @Produce 	json
// @Param       from           query string true "Start of range (from=2026-01-10T12:15:31+03:00)"
// @Param 		to             query string true "End of range   (to=2026-03-22T19:07:13+03:00)"
// @Param 		product-id     query string false "If you need a specific product"
// @Param 		category-id    query string false "If you need a specific category"
// @Param 		payment-method query string false "If you need a specific payment method"
// @Param 		group-by       query string false "Grouping by different fields (group-by=day/month/week/category/product)"
// @Param 		sort-by        query string false "Sorting by field (sort-by=group/value)"
// @Param 		sort-order     query string false "Sorting order (sort-order=asc/desc)"
// @Success 	200 {object} []response.AnalyticsResponse
// @Failure 	400 {object} response.Error
// @Failure 	500 {object} response.Error
// @Router 		/v1/analytics/p90 [get]
func (r *V1) getPercentile90(ctx *fiber.Ctx) error {
	params, s := parseAnalyticsParams(ctx)
	if s != "" {
		return errorResponse(ctx, fiber.StatusBadRequest, s)
	}

	results, err := r.a.GetPercentile90(ctx.UserContext(), params)
	if err != nil {
		r.l.Error(err, "restapi - v1 - getPercentile90")

		return errorResponse(ctx, fiber.StatusInternalServerError, "internal error")
	}

	resp := make([]response.AnalyticsResponse, len(results))
	for i, r := range results {
		resp[i] = response.AnalyticsResponse{
			GroupLabel: r.GroupLabel,
			Value:      r.Value,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func parseAnalyticsParams(ctx *fiber.Ctx) (entity.AnalyticsParams, string) {
	// From
	from, err := time.Parse(time.RFC3339, ctx.Query("from"))
	if err != nil {
		return entity.AnalyticsParams{}, "invalid from"
	}

	// To
	to, err := time.Parse(time.RFC3339, ctx.Query("to"))
	if err != nil {
		return entity.AnalyticsParams{}, "invalid to"
	}

	params := entity.AnalyticsParams{
		From: from,
		To:   to,
	}

	// Product ID
	if v := ctx.Query("product-id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return entity.AnalyticsParams{}, "invalid product-id"
		}
		params.ProductID = &id
	}

	// Category ID
	if v := ctx.Query("category-id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return entity.AnalyticsParams{}, "invalid category-id"
		}
		params.CategoryID = &id
	}

	if v := ctx.Query("payment-method"); v != "" {
		if v != "cash" && v != "card" && v != "online" {
			return entity.AnalyticsParams{}, "invalid payment-method: must be one of: \"cash\" / \"card\" / \"online\""
		}
		params.PaymentMethod = &v
	}

	if v := ctx.Query("group-by"); v != "" {
		gb := entity.GroupBy(v)
		switch gb {
		case entity.GroupByDay, entity.GroupByWeek, entity.GroupByMonth, entity.GroupByCategory, entity.GroupByProduct:
			params.GroupBy = &gb
		default:
			return entity.AnalyticsParams{}, "invalid group-by: must be on of: \"day\" / \"week\" / \"month\" / \"category\" / \"product\""
		}
	}

	if v := ctx.Query("sort-by"); v != "" {
		switch v {
		case "group", "value":
			params.SortBy = &v
		default:
			return entity.AnalyticsParams{}, "invalid sort-by: must be one of: \"group\" / \"value\""
		}
	}

	if v := ctx.Query("sort-order"); v != "" {
		so := entity.SortOrder(v)
		switch so {
		case entity.SortAscL, entity.SortDescL, entity.SortAscU, entity.SortDescU:
			params.SortOrder = so
		default:
			return entity.AnalyticsParams{}, "invalid sort-order: must be one of: \"asc\" / \"desc\" / \"ASC\" / \"DESC\""
		}
	}

	return params, ""
}
