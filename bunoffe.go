// Bunoffe is a small library to facilitate testing bun's (https://github.com/uptrace/bun)
// queries. One should feel free to copy and paste it directly into
// the code he/she is working on.
package bunoffe

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type (
	// Executor is the interface that wraps the methods of a query
	// executor type. Bun's queries can be executed with one of the
	// following methods: Exec, Scan, and Exists. Instead of calling
	// them directly, when using and Executor, you should use them
	// indirectly. For instance:
	//
	//     err := executor.Scan(
	//         ctx,
	//         db.NewSelect().Model(&m).WherePK(),
	//     )
	Executor interface {
		Exec(context.Context, ExecQuery, ...any) (sql.Result, error)
		Scan(context.Context, ScanQuery, ...any) error
		Exists(context.Context, ExistsQuery) (bool, error)
	}

	// ExecQuery is the interface that wraps the method Exec. Every
	// bun query can run Exec.
	//
	// Besides de Exec method, the GetModel method is required for
	// the MockQueryExecutor.
	ExecQuery interface {
		Exec(context.Context, ...any) (sql.Result, error)
		GetModel() bun.Model
	}

	// ScanQuery is the interface that wraps the method Scan.
	//
	// Besides de Exec method, the GetModel method is required for
	// the MockQueryExecutor.
	ScanQuery interface {
		Scan(context.Context, ...any) error
		GetModel() bun.Model
	}

	// ExistsQuery is the interface that wraps the method Exists.
	//
	// Besides de Exec method, the GetModel method is required for
	// the MockQueryExecutor.
	ExistsQuery interface {
		Exists(context.Context) (bool, error)
		GetModel() bun.Model
	}

	// QueryRealizer is the type of a Executor that executes the queries
	// that are passed to one of its methods. Using the realizer has the
	// same effect of executing a bun query directly.
	QueryRealizer struct{}
)

// Exec executes a bun query that has the Exec method. Calling:
//
//	executor.Exec(ctx, query, args...)
//
// is equivalent to running
//
//	query.Exec(ctx, args...)
func (QueryRealizer) Exec(
	ctx context.Context,
	q ExecQuery,
	args ...any,
) (sql.Result, error) {
	return q.Exec(ctx, args...)
}

// Scan executes a bun query that has the Scan method. Calling:
//
//	executor.Scan(ctx, query, args...)
//
// is equivalent to running
//
//	query.Scan(ctx, args...)
func (QueryRealizer) Scan(ctx context.Context, q ScanQuery, args ...any) error {
	return q.Scan(ctx, args...)
}

// Exists executes a bun query that has the Exists method. Calling:
//
//	executor.Exists(ctx, query)
//
// is equivalent to running
//
//	query.Exists(ctx)
func (QueryRealizer) Exists(ctx context.Context, q ExistsQuery) (bool, error) {
	return q.Exists(ctx)
}
