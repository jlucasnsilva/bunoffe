package bunoffe

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

type (
	// MockQueryExecutor is an Executor that ignores the queries
	// passed to its methods (Exec, Scan, and Exists). Instead,
	// the returned values and values assigned to the model are
	// the ones provided to operations (Ops field).
	MockQueryExecutor struct {
		// Ops is a slice of operations. Each time an Executor method
		// is called, next operation in line (starting with the first)
		// will be executed.
		Ops []MockedQueryOperation
		idx int
	}

	// MockedQueryOperation is interface that works as common type
	// for all mock operations.
	MockedQueryOperation interface {
		doNothing()
	}

	// MockExecOperation is a type to mock a Exec call.
	MockExecOperation struct {
		// If Model is not nil and Error is nil, when Exec is called, it will
		// contain the value passed to the query method `.Model(&m)`.
		Model any

		// If Args is not nil and Error is nil, when Exec is called, each of
		// its values will be assigned to parameter `...args`.
		Args []any

		// If Result is not nil and Error is nil, when Exec is called, it will
		// return Result.
		Result sql.Result

		// If Error is not nil, Exec will return a nil sql.Result and this
		// Error.
		Error error
	}

	// MockScanOperation is a type to mock a Scan call.
	MockScanOperation struct {
		// If Model is not nil and Error is nil, when Scan is called, it will
		// be assigned the value passed to the query method `.Model(&m)`.
		Model any

		// If Args is not nil and Error is nil, when Exec is called, each of
		// its values will be assigned to parameter `...args`.
		Args []any

		// If Error is not nil, Scan will return it.
		Error error
	}

	MockExistsOperation struct {
		// If Error is not nil, this value will be returned when Exists is
		// called. Otherwise false is returned.
		Exists bool

		// If Error is not nil, Scan will return it.
		Error error
	}

	MockQueryResult struct {
		LastInsertIdValue int64
		LastInsertIdError error

		RowsAffectedValue int64
		RowsAffectedError error
	}
)

func (MockExecOperation) doNothing()   {}
func (MockScanOperation) doNothing()   {}
func (MockExistsOperation) doNothing() {}

// Creates a *bun.DB with a mocked database.
func NewMockedBunDB() (*bun.DB, error) {
	sqldb, _, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	return bun.NewDB(sqldb, sqlitedialect.New()), nil
}

// Exec mocks a query.Exec call. See the MockExecOperation documentation for details.
func (ex *MockQueryExecutor) Exec(
	ctx context.Context,
	q ExecQuery,
	args ...any,
) (sql.Result, error) {
	nop := ex.nextOp()
	op, ok := nop.(MockExecOperation)
	if !ok {
		panic(opCastError("MockExec", nop))
	}

	if op.Error != nil {
		return nil, op.Error
	}

	if op.Model != nil {
		assign(
			reflect.ValueOf(op.Model),
			reflect.ValueOf(q.GetModel().Value()),
		)
	}

	if len(op.Args) > 0 && len(op.Args) != len(args) {
		panic("operation.Args and args should have the same length")
	}
	for i, val := range op.Args {
		assign(
			reflect.ValueOf(args[i]),
			reflect.ValueOf(val),
		)
	}
	return op.Result, nil
}

// Exec mocks a query.Scan call. See the MockScanOperation documentation for details.
func (ex *MockQueryExecutor) Scan(ctx context.Context, q ScanQuery, args ...any) error {
	nop := ex.nextOp()
	op, ok := nop.(MockScanOperation)
	if !ok {
		panic(opCastError("MockScan", nop))
	}

	if op.Error != nil {
		return op.Error
	}

	if op.Model != nil {
		assign(
			reflect.ValueOf(q.GetModel().Value()),
			reflect.ValueOf(op.Model),
		)
	}
	for i, val := range op.Args {
		assign(
			reflect.ValueOf(args[i]),
			reflect.ValueOf(val),
		)
	}
	return nil
}

// Exec mocks a query.Exists call. See the MockExistsOperation documentation for details.
func (ex *MockQueryExecutor) Exists(ctx context.Context, q ExistsQuery) (bool, error) {
	nop := ex.nextOp()
	op, ok := nop.(MockExistsOperation)
	if !ok {
		panic(opCastError("MockExists", nop))
	}

	if op.Error != nil {
		return false, op.Error
	}
	return op.Exists, nil
}

func (ex *MockQueryExecutor) nextOp() MockedQueryOperation {
	if len(ex.Ops) <= ex.idx {
		s := fmt.Sprintf(
			"mocked query requested operation #%v, but test only contains %v",
			ex.idx,
			len(ex.Ops),
		)
		panic(s)
	}

	ex.idx++
	return ex.Ops[ex.idx-1]
}

func (r MockQueryResult) LastInsertId() (int64, error) {
	return r.LastInsertIdValue, r.LastInsertIdError
}

func (r MockQueryResult) RowsAffected() (int64, error) {
	return r.RowsAffectedValue, r.RowsAffectedError
}

func opCastError(expected string, found any) string {
	return fmt.Sprintf("expected '%v' operation, but found '%T'", expected, found)
}

func assign(dest reflect.Value, src reflect.Value) {
	switch {
	case dest.Kind() == reflect.Ptr && src.Kind() == reflect.Ptr:
		dest.Elem().Set(src.Elem())
	case dest.Kind() == reflect.Ptr && src.Kind() != reflect.Ptr:
		dest.Elem().Set(src)
	case dest.Kind() != reflect.Ptr && src.Kind() == reflect.Ptr:
		dest.Set(src.Elem())
	case dest.Kind() != reflect.Ptr && src.Kind() != reflect.Ptr:
		dest.Set(src.Elem())
	}
}
