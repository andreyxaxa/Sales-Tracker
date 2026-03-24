package analytics

import (
	"context"
	"fmt"
	"strings"

	"github.com/andreyxaxa/Sales-Tracker/internal/entity"
	"github.com/andreyxaxa/Sales-Tracker/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type AnalyticsRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AnalyticsRepo {
	return &AnalyticsRepo{pg}
}

func (r *AnalyticsRepo) GetRevenue(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	groupSelect, groupByClause, joinClause, whereClause, args := r.buildBase(params)

	sql := fmt.Sprintf(`
		SELECT %s,
			COALESCE(SUM(s.quantity * s.unit_price), 0) AS value
		FROM sales s
		%s
		%s
		%s
		%s
	`, groupSelect, joinClause, whereClause, groupByClause, orderClause(params))

	res, err := r.collectRows(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsRepo - GetRevenue - r.collectRows: %w", err)
	}

	return res, nil
}

func (r *AnalyticsRepo) GetAvgRevenue(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	groupSelect, groupByClause, joinClause, whereClause, args := r.buildBase(params)

	sql := fmt.Sprintf(`
		SELECT %s,
			COALESCE(AVG(s.quantity * s.unit_price), 0) AS value
		FROM sales s
		%s
		%s
		%s
		%s
	`, groupSelect, joinClause, whereClause, groupByClause, orderClause(params))

	res, err := r.collectRows(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsRepo - GetAvgRevenue - r.collectRows: %w", err)
	}

	return res, nil
}

func (r *AnalyticsRepo) GetSalesCount(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	groupSelect, groupByClause, joinClause, whereClause, args := r.buildBase(params)

	sql := fmt.Sprintf(`
		SELECT %s,
			COALESCE(COUNT(*), 0) AS value
		FROM sales s
		%s
		%s
		%s
		%s
	`, groupSelect, joinClause, whereClause, groupByClause, orderClause(params))

	res, err := r.collectRows(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsRepo - GetSalesCount - r.collectRows: %w", err)
	}

	return res, nil
}

func (r *AnalyticsRepo) GetMedian(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	groupSelect, groupByClause, joinClause, whereClause, args := r.buildBase(params)

	sql := fmt.Sprintf(`
		SELECT %s,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY s.quantity * s.unit_price), 0) AS value
		FROM sales s
		%s
		%s
		%s
		%s
	`, groupSelect, joinClause, whereClause, groupByClause, orderClause(params))

	res, err := r.collectRows(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsRepo - GetMedian - r.collectRows: %w", err)
	}

	return res, nil
}

func (r *AnalyticsRepo) GetPercentile90(ctx context.Context, params entity.AnalyticsParams) ([]entity.AnalyticsResult, error) {
	groupSelect, groupByClause, joinClause, whereClause, args := r.buildBase(params)

	sql := fmt.Sprintf(`
		SELECT %s,
			COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY s.quantity * s.unit_price), 0) AS value
		FROM sales s
		%s
		%s
		%s
		%s
	`, groupSelect, joinClause, whereClause, groupByClause, orderClause(params))

	res, err := r.collectRows(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsRepo - GetPercentile90 - r.collectRows: %w", err)
	}

	return res, nil
}

func (r *AnalyticsRepo) buildBase(params entity.AnalyticsParams) (
	groupSelect string,
	groupByClause string,
	joinClause string,
	whereClause string,
	args []any,
) {
	argIdx := 1
	var conditions []string

	productJoin := false
	categoryJoin := false

	// From
	conditions = append(conditions, fmt.Sprintf("s.sold_at >= $%d", argIdx))
	args = append(args, params.From)
	argIdx++

	// To
	conditions = append(conditions, fmt.Sprintf("s.sold_at <= $%d", argIdx))
	args = append(args, params.To)
	argIdx++

	// where
	if params.ProductID != nil {
		conditions = append(conditions, fmt.Sprintf("s.product_id = $%d", argIdx))
		args = append(args, params.ProductID)
		argIdx++
	} else if params.CategoryID != nil {
		productJoin = true
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argIdx))
		args = append(args, params.CategoryID)
		argIdx++
	}

	if params.PaymentMethod != nil {
		conditions = append(conditions, fmt.Sprintf("s.payment_method = $%d", argIdx))
		args = append(args, params.PaymentMethod)
	}

	whereClause = "WHERE " + strings.Join(conditions, " AND ")

	// group by
	switch {
	case params.GroupBy == nil:
		groupSelect = "'total' AS group_label"
		groupByClause = ""
	case *params.GroupBy == entity.GroupByDay:
		groupSelect = "TO_CHAR(DATE_TRUNC('day', s.sold_at), 'YYYY-MM-DD') AS group_label"
		groupByClause = "GROUP BY DATE_TRUNC('day', s.sold_at)"
	case *params.GroupBy == entity.GroupByWeek:
		groupSelect = "TO_CHAR(DATE_TRUNC('week', s.sold_at), 'YYYY-MM-DD') AS group_label"
		groupByClause = "GROUP BY DATE_TRUNC('week', s.sold_at)"
	case *params.GroupBy == entity.GroupByMonth:
		groupSelect = "TO_CHAR(DATE_TRUNC('month', s.sold_at), 'YYYY-MM') AS group_label"
		groupByClause = "GROUP BY DATE_TRUNC('month', s.sold_at)"
	case *params.GroupBy == entity.GroupByCategory:
		categoryJoin = true
		groupSelect = "c.name AS group_label"
		groupByClause = "GROUP BY c.name"
	case *params.GroupBy == entity.GroupByProduct:
		productJoin = true
		groupSelect = "p.name AS group_label"
		groupByClause = "GROUP BY p.name"
	default:
		groupSelect = "'total' AS group_label"
		groupByClause = ""
	}

	// join
	if categoryJoin {
		joinClause = "JOIN products p ON s.product_id = p.id JOIN categories c ON p.category_id = c.id"
	} else if productJoin {
		joinClause = "JOIN products p ON s.product_id = p.id"
	}

	return groupSelect, groupByClause, joinClause, whereClause, args
}

func orderClause(params entity.AnalyticsParams) string {
	allowed := map[string]string{
		"group": "group_label",
		"value": "value",
	}

	sortBy := "group_label"
	if params.SortBy != nil {
		if col, ok := allowed[*params.SortBy]; ok {
			sortBy = col
		}
	}

	order := "ASC"
	if params.SortOrder == entity.SortDescU || params.SortOrder == entity.SortDescL {
		order = "DESC"
	}

	return fmt.Sprintf("ORDER BY %s %s", sortBy, order)
}

func (r *AnalyticsRepo) collectRows(ctx context.Context, sql string, args []any) ([]entity.AnalyticsResult, error) {
	executor := r.GetExecutor(ctx)

	rows, err := executor.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AnalyticsRepo - collectRows - executor.Query: %w", err)
	}
	defer rows.Close()

	res, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.AnalyticsResult])
	if err != nil {
		return nil, fmt.Errorf("AnalyticsRepo - collectRows - pgx.CollectRows: %w", err)
	}

	return res, nil
}
