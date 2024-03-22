# Overview

Bunoffe is a small library to facilitate testing [bun's](https://github.com/uptrace/bun)
queries. One should feel free to copy and paste it directly into
the code he/she is working on.

# Usage

## Exec

Instead of writing

```go
db := bun.NewDB(sqldb, sqlitedialect.New())

result, err := bundb.NewExec().
    Model(&m).
    Exec(ctx)
```

Do

```go
db := bun.NewDB(sqldb, sqlitedialect.New())
executor := bunoffe.QueryRealizer{}

err := executor.Exec(
    ctx,
    bundb.NewExec().
        Model(&m),
)
```

## Exists

Instead of writing

```go
db := bun.NewDB(sqldb, sqlitedialect.New())

exists, err := bundb.NewSelect().
    Model(&m).
    Exists(ctx)
```

Do

```go
db := bun.NewDB(sqldb, sqlitedialect.New())
executor := bunoffe.QueryRealizer{}

exists, err := executor.Exists(
    ctx,
    bundb.NewSelect().
        Model(&m),
)
```

## Scan

Instead of writing

```go
db := bun.NewDB(sqldb, sqlitedialect.New())

err := bundb.NewSelect().
    Model(&m).
    Scan(ctx)
```

Do

```go
db := bun.NewDB(sqldb, sqlitedialect.New())
executor := bunoffe.QueryRealizer{}

err := executor.Scan(
    ctx,
    bundb.NewSelect().
        Model(&m),
)
```

# Testing

Bunoffe provides a set mocked operations. Check it out.
