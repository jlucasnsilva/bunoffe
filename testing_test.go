package bunoffe

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type model struct {
	String string
	Int    int
}

func TestMocks(t *testing.T) {
	db, err := NewMockedBunDB()
	require.Nil(t, err)

	ctx := context.Background()

	t.Run("test exec", func(t *testing.T) {
		// expected
		var (
			err     = errors.New("an error")
			m       = model{String: "Hello, world!", Int: 33}
			message = "hadouken"
			pi      = 3.14
			result  = MockQueryResult{
				LastInsertIdValue: 10,
				RowsAffectedValue: 11,
			}
		)

		ex := MockQueryExecutor{
			Ops: []MockedQueryOperation{
				MockScanOperation{},
				MockExecOperation{Error: err},
				MockExecOperation{Result: result, Model: &m},
				MockExecOperation{Args: []any{message, pi}},
			},
		}

		// results
		var (
			e error
			n model
			s string
			f float64
			r sql.Result
		)

		assert.Panics(t, func() {
			ex.Exec(
				ctx,
				db.NewInsert().Model(&n),
			)
		})

		n = model{}
		r, e = ex.Exec(
			ctx,
			db.NewInsert().Model(&n),
		)
		assert.NotNil(t, e)
		assert.Nil(t, r)

		n = model{}
		r, e = ex.Exec(
			ctx,
			db.NewInsert().Model(&n),
		)
		assert.Nil(t, e)
		assert.Equal(t, result, r)
		assert.Equal(t, m, n)

		n = model{}
		r, e = ex.Exec(
			ctx,
			db.NewInsert().Model(&n),
			&s, &f,
		)
		assert.Nil(t, e)
		assert.Nil(t, r)
		assert.Equal(t, m, n)
		assert.Equal(t, message, s)
		assert.Equal(t, pi, f)
	})

	t.Run("test scan", func(t *testing.T) {
		// expected
		var (
			err     = errors.New("an error")
			m       = model{String: "Hello, world!", Int: 33}
			message = "hadouken"
			pi      = 3.14
		)

		ex := MockQueryExecutor{
			Ops: []MockedQueryOperation{
				MockExecOperation{},
				MockScanOperation{Error: err},
				MockScanOperation{Model: &m},
				MockScanOperation{Model: &m, Args: []any{message, pi}},
			},
		}

		// results
		var (
			e error
			n model
			s string
			f float64
		)

		assert.Panics(t, func() {
			ex.Scan(
				ctx,
				db.NewSelect().Model(&n),
			)
		})

		n = model{}
		e = ex.Scan(
			ctx,
			db.NewSelect().Model(&n),
		)
		assert.NotNil(t, e)

		n = model{}
		e = ex.Scan(
			ctx,
			db.NewSelect().Model(&n),
		)
		assert.Nil(t, e)
		assert.Equal(t, m, n)

		n = model{}
		e = ex.Scan(
			ctx,
			db.NewSelect().Model(&n),
			&s, &f,
		)
		assert.Nil(t, e)
		assert.Equal(t, m, n)
		assert.Equal(t, message, s)
		assert.Equal(t, pi, f)
	})

	t.Run("test exists", func(t *testing.T) {
		// expected
		err := errors.New("an error")
		ex := MockQueryExecutor{
			Ops: []MockedQueryOperation{
				MockExecOperation{},
				MockExistsOperation{Error: err},
				MockExistsOperation{Exists: true},
				MockExistsOperation{Exists: false},
			},
		}

		// results
		var (
			n model
			e error
			f bool
		)

		assert.Panics(t, func() {
			ex.Exists(
				ctx,
				db.NewSelect().Model(&n),
			)
		})

		n = model{}
		f, e = ex.Exists(
			ctx,
			db.NewSelect().Model(&n),
		)
		assert.False(t, f)
		assert.NotNil(t, e)

		n = model{}
		f, e = ex.Exists(
			ctx,
			db.NewSelect().Model(&n),
		)
		assert.Nil(t, e)
		assert.True(t, f)

		n = model{}
		f, e = ex.Exists(
			ctx,
			db.NewSelect().Model(&n),
		)
		assert.Nil(t, e)
		assert.False(t, f)
	})
}
