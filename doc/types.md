# nice-expr Types and Operations

## Types

### Primitive Types

* `int`: 64-bit integer
* `dec`: 64-bit float
* `str`: string
* `bool`: boolean

### Compound Types

* `list[T]`: list
* `map[K]V`: map

## Zero-Values

Each type has a zero value when declared on a variable with no explicit assignment.

| Type      | Zero Value |
| --------- | ---------- |
| `int`     | `0`        |
| `dec`     | `0.0`      |
| `str`     | `""`       |
| `bool`    | `false`    |
| `list[T]` | `[]`       |
| `map[K]V` | `<\|\|>`   |

## Operations

All operations' outputs match their left value type.

### Unary

| Op    | `int`  | `dec`  | `str` | `bool` | `list[T]` | `map[K]V` |
| ----- | ------ | ------ | ----- | ------ | --------- | --------- |
| `not` |        |        |       | Invert |           |           |
| `-`   | Negate | Negate |       |        |           |           |

### Binary Arithmetic

| Op  | `int`          | `dec`          | `str`         | `bool` | `list[T]`     | `map[K]V` |
| --- | -------------- | -------------- | ------------- | ------ | ------------- | --------- |
| `+` | Addition       | Addition       | Concatenation |        | Concatenation |           |
| `-` | Subtraction    | Subtraction    | Intersection  |        | Intersection  |           |
| `*` | Multiplication | Multiplication |               |        |               |           |
| `/` | Division       | Division       |               |        |               |           |
| `%` | Modulo         |                |               |        |               |           |

### Binary Comparison

Both sides of the expression must match in type.
Always returns a boolean.

| Op   | `int`    | `dec`    | `str`           | `bool`   | `list[T]`       | `map[K]V` |
| ---- | -------- | -------- | --------------- | -------- | --------------- | --------- |
| `=`  | Equality | Equality | Equality        | Equality | Equality        | Equality  |
| `>`  | Gt       | Gt       | Lexicographical |          | Lexicographical |           |
| `>=` | Ge       | Ge       |                 |          |                 |           |
| `<`  | Lt       | Lt       |                 |          |                 |           |
| `<=` | Le       | Le       |                 |          |                 |           |

### Binary Logical

Both sides must be a boolean.
Always returns a boolean.

* `and`
* `or`

### Indexing Operator `collection_indexer`

Only operates on `str`, `list`, and `map` on the left.
Type of the indexer differs based on the type.
The return type differs based on the type.

Returns the zero value if the indexer is not in the collection.

| Type      | Indexing Type | Return Type | Returns                                   |
| --------- | ------------- | ----------- | ----------------------------------------- |
| `str`     | `int`         | `str`       | the len-1 string at that index, 0-indexed |
| `list[T]` | `int`         | `T`         | the element at that index, 0-indexed      |
| `map[K]V` | `K`           | `V`         | the value associated with the key         |
