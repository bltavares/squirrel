// +build go1.8

package squirrel

import (
	"context"
	"database/sql"
)

// Execer is the interface that wraps the Exec method.
//
// Exec executes the given query as implemented by database/sql.Exec.
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Queryer is the interface that wraps the Query method.
//
// Query executes the given query as implemented by database/sql.Query.
type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// QueryRower is the interface that wraps the QueryRow method.
//
// QueryRow executes the given query as implemented by database/sql.QueryRow.
type QueryRower interface {
	QueryRow(query string, args ...interface{}) RowScanner
}

type DbRunnerWithContext struct {
	ctx    context.Context
	runner *dbRunner
}

func (r *dbRunner) WithContext(ctx context.Context) DbRunnerWithContext {
	DbRunnerWithContext{
		ctx:    ctx,
		runner: r,
	}
}

func (r DbRunnerWithContext) QueryRow(query string, args ...interface{}) RowScanner {
	return r.runner.DB.WithContext(r.ctx).QueryRow(query, args...)
}

type TxRunnerWithContext struct {
	ctx    context.Context
	runner *txRunner
}

func (r *txRunner) WithContext(ctx context.Context) TxRunnerWithContext {
	TxRunnerWithContext{
		ctx:    ctx,
		runner: r,
	}
}

func (r *txRunner) QueryRow(query string, args ...interface{}) RowScanner {
	return r.runner.Tx.WithContext(r.ctx).QueryRow(query, args...)
}

// ExecContextWith ExecContexts the SQL returned by s with db.
func ExecContextWith(ctx context.Context, db Execer, s Sqlizer) (res sql.Result, err error) {
	query, args, err := s.ToSql()
	if err != nil {
		return
	}
	return db.WithContext(ctx).Exec(ctx, query, args...)
}

// QueryContextWith QueryContexts the SQL returned by s with db.
func QueryContextWith(ctx context.Context, db Queryer, s Sqlizer) (rows *sql.Rows, err error) {
	query, args, err := s.ToSql()
	if err != nil {
		return
	}
	return db.WithContext(ctx).Query(query, args...)
}

// QueryRowContextWith QueryRowContexts the SQL returned by s with db.
func QueryRowContextWith(ctx context.Context, db QueryRower, s Sqlizer) RowScanner {
	query, args, err := s.ToSql()
	return &Row{RowScanner: db.WithContext(ctx).QueryRow(query, args...), err: err}
}
