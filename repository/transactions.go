package repository

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type DBConn interface {
	bun.IDB
	bun.IConn
}

var _ DBConn = (*bun.DB)(nil)
var _ DBConn = bun.Tx{}

type DBWrapper struct {
	DBConn
}

func NewDBWrapper(db *bun.DB) DBWrapper {
	return DBWrapper{DBConn: db}
}

func (w DBWrapper) WithTx(tx bun.Tx) DB {
	return &DBWrapper{DBConn: tx}
}

type DB interface {
	DBConn
	WithTx(tx bun.Tx) DB
}

var _ DB = DBWrapper{}

func RunInTx(ctx context.Context, db *bun.DB, fn func(tx bun.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
