// +build go1.8

package squirrel

import (
	"context"
	"database/sql"

	"github.com/lann/builder"
)

type SelectDataWithContext struct {
	ctx  context.Context
	data *selectData
}

func (d *selectData) WithContext(ctx context.Context) SelectDataWithContext {
	SelectDataWithContext{
		ctx:  ctx,
		data: d,
	}
}

func (d SelectDataWithContext) Exec() (sql.Result, error) {
	if d.data.RunWith == nil {
		return nil, RunnerNotSet
	}
	return ExecContextWith(d.ctx, d.data.RunWith, d.data)
}

func (d SelectDataWithContext) Query() (*sql.Rows, error) {
	if d.data.RunWith == nil {
		return nil, RunnerNotSet
	}
	return QueryContextWith(d.ctx, d.data.RunWith, d.data)
}

func (d SelectDataWithContext) QueryRow() RowScanner {
	if d.data.RunWith == nil {
		return &Row{err: RunnerNotSet}
	}
	queryRower, ok := d.data.RunWith.(QueryRower)
	if !ok {
		return &Row{err: RunnerNotQueryRunner}
	}
	return QueryRowContextWith(d.ctx, queryRower, d.data)
}

type SelectBuilderWithContext struct {
	ctx     context.Context
	builder SelectBuilder
}

func (b SelectBuilder) WithContext(ctx context.Context) SelectBuilderWithContext {
	SelectBuilderWithContext{
		ctx:     ctx,
		builder: b,
	}
}

func (b SelectBuilderWithContext) Exec() (sql.Result, error) {
	data := builder.GetStruct(b.builder).(selectData)
	return data.WithContext(b.ctx).Exec()
}

func (b SelectBuilderWithContext) Query() (*sql.Rows, error) {
	data := builder.GetStruct(b.builder).(selectData)
	return data.WithContext(b.ctx).Query()
}

func (b SelectBuilderWithContext) QueryRow() RowScanner {
	data := builder.GetStruct(b.builder).(selectData)
	return data.WithContext(b.ctx).QueryRow()
}
