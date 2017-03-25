// +build go1.8

package squirrel

import (
	"context"
	"database/sql"
)

type DbRunnerWithContext struct {
	ctx    context.Context
	runner *dbRunner
}

func (r *dbRunner) WithContext(ctx context.Context) DbRunnerWithContext {
	return DbRunnerWithContext{
		ctx:    ctx,
		runner: r,
	}
}

func (r DbRunnerWithContext) QueryRow(query string, args ...interface{}) RowScanner {
	return r.runner.DB.QueryRowContext(r.ctx, query, args...)
}

type TxRunnerWithContext struct {
	ctx    context.Context
	runner *txRunner
}

func (r *txRunner) WithContext(ctx context.Context) TxRunnerWithContext {
	return TxRunnerWithContext{
		ctx:    ctx,
		runner: r,
	}
}

func (r TxRunnerWithContext) QueryRow(ctx context.Context, query string, args ...interface{}) RowScanner {
	return r.runner.Tx.QueryRowContext(r.ctx, query, args...)

}

type ContextAwareExec struct {
	ctx    context.Context
	execer Execer
}

func (e Execer) WithContext(ctx context.Context) ContextAwareExec {
	return ContextAwareExec{
		ctx:    ctx,
		execer: e,
	}
}

func (e ContextAwareExec) Exec() (res sql.Result, err error) {
	return nil, nil
}

// ExecContextWith ExecContexts the SQL returned by s with db.
func ExecContextWith(ctx context.Context, db Execer, s Sqlizer) (res sql.Result, err error) {
	query, args, err := s.ToSql()
	if err != nil {
		return
	}
	return db.WithContext(ctx).Exec(query, args...)
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
