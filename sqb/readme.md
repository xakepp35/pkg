# QueryBuilder Usage Guide

## Introduction
`Sqb` is a lightweight SQL query builder.

## Installation
To use `SQB`, include the `sqb` package in your Go project.

```go
import "github.com/xakepp35/pkg/sqb"
```

## Quick Start
```go
func main() {
	// Acquiring a QueryBuilder Instance
    qb := sqb.GetBuilder()
    defer sqb.RealiseBuilder(qb) // Release the builder after use

    // Создаем запрос
    query, args := qb.Select("id", "name").
        From("users").
        Where("age > ?", 18).
        Where("status = ?", "active").
        Limit(10).
        Offset(5).
        Build()

    // Output result for debug
    fmt.Println("Query:", query) // Output: SELECT id, name FROM users WHERE (age > $1) AND (status = $2) LIMIT 10 OFFSET 5
    fmt.Println("Args:", args)   // Output: [18 "active"]
}
```

## Getting Started
### Acquiring a QueryBuilder Instance
To start building a query, obtain a `SQB` instance from the pool:

```go
qb := sqb.GetBuilder()
defer sqb.RealiseBuilder(qb) // Ensure the builder is released after use
```

## Building Queries
### Select Statement
```go
query, args := sqb.GetBuilder().Select("id", "name").From("users").Build()
fmt.Println(query) // Output: "SELECT id, name FROM users"
fmt.Println(args)  // Output: []
```

### Where Clause with Parameters
```go
query, args := sqb.GetBuilder().
    Select("id", "name").
    From("users").
    Where("age > ?", 18).
    Build()
fmt.Println(query) // Output: "SELECT id, name FROM users WHERE (age > $1)"
fmt.Println(args)  // Output: [18]
```

### Using Multiple Where Conditions
```go
query, args := sqb.GetBuilder().
    Select("id", "name").
    From("users").
    Where("age > ?", 18).
    Where("status = ?", "active").
    Build()
fmt.Println(query) // Output: "SELECT id, name FROM users WHERE (age > $1) AND (status = $2)"
fmt.Println(args)  // Output: [18 "active"]
```

### Using OR Operator
```go
query, args := sqb.GetBuilder().
    Select("id", "name").
    From("users").
    Where("age > ?", 18).
    Or().
	Where("status = ?", "active").
    Build()
fmt.Println(query) // Output: "SELECT id, name FROM users WHERE (age > $1) OR (status = $2)"
fmt.Println(args)  // Output: [18 "active"]
```

### Adding Limit and Offset
```go
query, args := sqb.GetBuilder().
    Select("id", "name").
    From("users").
    Limit(10).
    Offset(5).
    Build()
fmt.Println(query) // Output: "SELECT id, name FROM users LIMIT 10 OFFSET 5"
fmt.Println(args)  // Output: []
```

## Resetting and Releasing the Builder
To reuse the builder, reset it:
```go
qb.Reset()
```

Once done, release the builder back to the pool:
```go
sqb.RealiseBuilder(qb)
```

## Summary
- Use `Build()` to generate the final SQL query and arguments.
- Always release the builder using `RealiseBuilder(qb)` after use.



