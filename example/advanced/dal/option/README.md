# DAL Options

This package provides utilities for defining and working with database query options using the `gorm` library. It includes abstractions for columns and entities, enabling dynamic query generation and improved code organization.

## Architecture Rules

### Column Methods (Single-Table Queries)

Column methods operate on individual columns and generate single-column query conditions. These are available on both `Column` and specific column types (e.g., `userColumn`, `orderColumn`).

**Examples:**

- `Eq(value)` - Equality condition
- `Gt(value)` - Greater than condition
- `Gte(value)` - Greater than or equal to condition
- `Lt(value)` - Less than condition
- `Lte(value)` - Less than or equal to condition

### Entity Methods (Multi-Column Operations)

Entity methods operate on entities and handle operations involving multiple columns. These are available on entity types (e.g., `userEntity`, `orderEntity`).

**Examples:**

- `Select(cols...)` - Select specific columns
- `Where(opts...)` - Apply multiple query conditions
- `OrderBy(columns...)` - Order by multiple columns
- `Limit(n)` - Limit result count
- `Offset(n)` - Offset results for pagination

## Files

- **`column.go`**: Defines the `Column` struct and methods for building single-column query conditions.
- **`order_option.go`**: Implements the `orderEntity` struct for managing order-related database operations. Provides entity-level methods for multi-column queries.
- **`user_option.go`**: Implements the `userEntity` struct for managing user-related database operations. Provides entity-level methods for multi-column queries.

## Usage

### Column-Level Queries (Single Column)

Use column methods for single-column conditions:

```go
// Single column condition
db.Where(USERS.Email.Eq("john@example.com"))
db.Where(ORDERS.Quantity.Gt(10))
```

### Entity-Level Queries (Multiple Columns)

Use entity methods for operations involving multiple columns:

```go
// Select specific columns
db = USERS.Select(USERS.ID, USERS.Email)

// Apply multiple conditions
db = USERS.Where(
    USERS.Email.Eq("john@example.com"),
    USERS.Name.Eq("John"),
)

// Order by multiple columns with limit and offset
db = ORDERS.OrderBy("created_at DESC", "id ASC")
db = ORDERS.Limit(10)
db = ORDERS.Offset(20)
```

### Combined Usage

You can combine column and entity methods for complex queries:

```go
db = USERS.Select(USERS.ID, USERS.Name, USERS.Email).
    Where(
        USERS.Email.Eq("john@example.com"),
        USERS.ID.Gt(100),
    ).
    OrderBy("created_at DESC").
    Limit(10)
```

## Dependencies

- [GORM](https://gorm.io/): An ORM library for Go.
- [gormqs](https://github.com/foxie-io/gormqs): A utility library for GORM.

## License

This project is licensed under the terms of the MIT license.
