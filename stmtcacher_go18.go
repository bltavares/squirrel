// +build go1.8

package squirrel

import (
	"context"
	"database/sql"
)

type StmtCacherWithContext struct {
	ctx    context.Context
	cacher *stmtCacher
}

func (sc *stmtCacher) WithContext(ctx context.Context) StmtCacherWithContext {
	return StmtCacherWithContext{
		ctx:    ctx,
		cacher: sc,
	}
}

func (sc StmtCacherWithContext) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	stmt, err := sc.cacher.Prepare(query)
	if err != nil {
		return
	}

	return stmt.WithContext(sc.ctx).Exec(args...)
}

func (sc StmtCacherWithContext) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	stmt, err := sc.cacher.Prepare(query)
	if err != nil {
		return
	}
	return stmt.WithContext(sc.ctx).Query(ctx, args...)
}

func (sc StmtCacherWithContext) QueryRow(query string, args ...interface{}) RowScanner {
	stmt, err := sc.cacher.Prepare(query)
	if err != nil {
		return &Row{err: err}
	}
	return stmt.WithContext(sc.ctx).QueryRow(args...)
}
