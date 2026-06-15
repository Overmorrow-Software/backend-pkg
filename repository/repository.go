package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Overmorrow-Software/backend-pkg/apierror"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository[M any, ID uuid.UUID | string | uint64 | int64] interface {
	Create(ctx context.Context, model *M) error
	GetByID(ctx context.Context, id ID) (*M, error)
	Update(ctx context.Context, model *M) error
	Delete(ctx context.Context, id ID) error
	SoftDelete(ctx context.Context, id ID) error
	Exists(ctx context.Context, model *M) (bool, error)
	FindAll(ctx context.Context) ([]M, error)
	FindAllWithOptions(ctx context.Context, opts *Options) ([]M, error)
	CountWithOptions(ctx context.Context, opts *Options) (int, error)
	FindPage(ctx context.Context, opts *Options) (PageResult[M], error)
	WithTx(tx bun.Tx) Repository[M, ID]
}

type PageResult[M any] struct {
	Items []M
	Total int
}

type Generic[M any, IDType uuid.UUID | string | uint64 | int64] struct {
	DB DB
}

var _ Repository[struct{}, uint64] = Generic[struct{}, uint64]{}

func NewGenericRepository[M any, IDType uuid.UUID | string | uint64 | int64](db *bun.DB) Generic[M, IDType] {
	return Generic[M, IDType]{DB: DBWrapper{DBConn: db}}
}

func (r Generic[M, IDType]) WithTx(tx bun.Tx) Repository[M, IDType] {
	return Generic[M, IDType]{DB: r.DB.WithTx(tx)}
}

func (r Generic[M, IDType]) Create(ctx context.Context, model *M) error {
	_, err := r.DB.NewInsert().Model(model).Exec(ctx)
	return err
}

func (r Generic[M, IDType]) GetByID(ctx context.Context, id IDType) (*M, error) {
	model := new(M)
	err := r.DB.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apierror.NotFound("not found")
	}
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r Generic[M, IDType]) FindAll(ctx context.Context) ([]M, error) {
	models := make([]M, 0)
	err := r.DB.NewSelect().Model(&models).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r Generic[M, IDType]) Update(ctx context.Context, model *M) error {
	_, err := r.DB.NewUpdate().Model(model).OmitZero().WherePK().Exec(ctx)
	return err
}

func (r Generic[M, IDType]) Delete(ctx context.Context, id IDType) error {
	_, err := r.DB.NewDelete().Model((*M)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r Generic[M, IDType]) SoftDelete(ctx context.Context, id IDType) error {
	_, err := r.DB.NewUpdate().Model((*M)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r Generic[M, IDType]) Exists(ctx context.Context, model *M) (bool, error) {
	return r.DB.NewSelect().Model(model).WherePK().Exists(ctx)
}

func (r Generic[M, IDType]) FindAllWithOptions(ctx context.Context, options *Options) ([]M, error) {
	models := make([]M, 0)
	q := r.DB.NewSelect().Model(&models)
	q = options.Apply(q)
	err := q.Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r Generic[M, IDType]) CountWithOptions(ctx context.Context, options *Options) (int, error) {
	q := r.DB.NewSelect().Model((*M)(nil))
	if options != nil {
		q = options.ApplyFilters(q)
		q = options.ApplyOrderBy(q)
	}
	return q.Count(ctx)
}

func (r Generic[M, IDType]) FindPage(ctx context.Context, opts *Options) (PageResult[M], error) {
	items, err := r.FindAllWithOptions(ctx, opts)
	if err != nil {
		return PageResult[M]{}, err
	}
	total, err := r.CountWithOptions(ctx, opts)
	if err != nil {
		return PageResult[M]{}, err
	}
	return PageResult[M]{Items: items, Total: total}, nil
}

func (o *Options) Apply(q *bun.SelectQuery) *bun.SelectQuery {
	if o != nil {
		q = o.ApplyFilters(q)
		q = o.ApplyOrderBy(q)
		q = o.ApplyPagination(q)
	}
	return q
}

func (o *Options) ApplyPagination(q *bun.SelectQuery) *bun.SelectQuery {
	if o == nil || q == nil {
		return q
	}
	if o.Pagination.IsValid() {
		q = q.Offset(int((o.Pagination.PageNum - 1) * o.Pagination.PageSize)) //nolint:gosec
		q = q.Limit(int(o.Pagination.PageSize))                               //nolint:gosec
	}
	return q
}

func (o *Options) ApplyOrderBy(q *bun.SelectQuery) *bun.SelectQuery {
	if o == nil || q == nil {
		return q
	}
	if o.Order.IsValid() {
		for col := range strings.SplitSeq(o.Order.OrderBy, ",") {
			col = strings.TrimSpace(col)
			if col == "" {
				continue
			}
			q = q.OrderExpr("? ?", bun.Ident(col), bun.Safe(o.Order.OrderType))
		}
	}
	return q
}

func (o *Options) ApplyFilters(q *bun.SelectQuery) *bun.SelectQuery {
	if o == nil || q == nil {
		return q
	}

	for _, filter := range o.Filters {
		filter.Operator = strings.ToLower(filter.Operator)
		if !filter.isValid() {
			continue
		}

		if filter.Operator == "in" {
			vals := make([]any, len(filter.Values))
			for i, v := range filter.Values {
				vals[i] = v
			}
			if filter.WhereOr {
				q = q.WhereOr("? IN (?)", bun.Ident(filter.Column), bun.List(vals))
			} else {
				q = q.Where("? IN (?)", bun.Ident(filter.Column), bun.List(vals))
			}
			continue
		}

		if filter.Operator == "like" || filter.Operator == "ilike" {
			filter.Value = fmt.Sprintf("%%%s%%", filter.Value)
		}

		if filter.Operator == "is" || filter.Operator == "is not" {
			if strings.ToLower(filter.Value) == "null" {
				q = q.Where("? ? NULL", bun.Ident(filter.Column), bun.Safe(filter.Operator))
				continue
			}
			q = q.Where("? ? ?", bun.Ident(filter.Column), bun.Safe(filter.Operator), bun.Safe(filter.Value))
			continue
		}

		if filter.WhereOr {
			q = q.WhereOr("? ? ?", bun.Ident(filter.Column), bun.Safe(filter.Operator), filter.Value)
		} else {
			q = q.Where("? ? ?", bun.Ident(filter.Column), bun.Safe(filter.Operator), filter.Value)
		}
	}

	return q
}
