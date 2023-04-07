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

### Casting

Cast a value to a type using `val as T`, where `val` is a value and `T` is a type.

| ↓ `val`; → `T` | `int` | `dec` | `str` | `bool` | `list[T]` | `map[K]V` |
| -------------- | ----- | ----- | ----- | ------ | --------- | --------- |
| `int`          | ✅     | ✅     | ✅     | ✅      |           |           |
| `dec`          | ✅     | ✅     | ✅     | ✅      |           |           |
| `str`          | ✅     | ✅     | ✅     | ✅      |           |           |
| `bool`         | ✅     | ✅     | ✅     | ✅      |           |           |
| `list[T]`      |       |       |       |        | ✅         |           |
| `map[K]V`      |       |       |       |        |           | ✅         |

When casting `list` and `map`, it will attempt to cast all the contained elements to the new type.

Trying a combination not in the above table is a runtime error.

When casting anything to bool, the returned value is whether or not the value is equal to to the zero value of its type.

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

### Precedence

From highest (first executed) to lowest (last executed).
Operators on the same row have the same precedence and are grouped left-to-right.

| Op Kind        | Operator      |
| -------------- | ------------- |
| Parens         | `()`          |
| Type Cast      | `as`          |
| Indexing       | `_`           |
| Unary Minus    | `-a`          |
| Multiplication | `* / %`       |
| Addition       | `+ -`         |
| Comparisons    | `< > <= >= =` |
| Logical        | `and or`      |
| Logical Not    | `not a`       |

#### BNF

* <https://docs.python.org/3/reference/expressions.html>
* <https://go.dev/ref/spec#Expressions>

```ebnf
Expr ::= IndexingExpr | "(" Expr ")";
IndexingExpr ::= UnaryMinusExpr "_" Expr;
```

### Unary

| Op    | `int`  | `dec`  | `str` | `bool` | `list[T]` | `map[K]V` |
| ----- | ------ | ------ | ----- | ------ | --------- | --------- |
| `not` |        |        |       | Invert |           |           |
| `-`   | Negate | Negate |       |        |           |           |

### Binary Arithmetic

| Op  | `int`          | `dec`          | `str`          | `bool` | `list[T]`      | `map[K]V` |
| --- | -------------- | -------------- | -------------- | ------ | -------------- | --------- |
| `+` | Addition       | Addition       | Concatenation  |        | Concatenation  |           |
| `-` | Subtraction    | Subtraction    | Set Difference |        | Set Difference |           |
| `*` | Multiplication | Multiplication |                |        |                |           |
| `/` | Division       | Division       |                |        |                |           |
| `%` | Modulo         |                |                |        |                |           |

The subtraction operator of `str` and `list` performs set difference on them,
removing from the left side elements in common with the right side.

### Binary Comparison

Both sides of the expression must match in type.
Always returns a boolean.

| Op   | `int`    | `dec`    | `str`           | `bool`   | `list[T]` | `map[K]V` |
| ---- | -------- | -------- | --------------- | -------- | --------- | --------- |
| `=`  | Equality | Equality | Equality        | Equality | Equality  | Equality  |
| `>`  | Gt       | Gt       | Lexicographical |          |           |           |
| `>=` | Ge       | Ge       | Lexicographical |          |           |           |
| `<`  | Lt       | Lt       | Lexicographical |          |           |           |
| `<=` | Le       | Le       | Lexicographical |          |           |           |

### Binary Logical

Both sides must be a boolean.
Always returns a boolean.

* `and`
* `or`

### Indexing Operator `collection_indexer`

Only operates on `str`, `list`, and `map` on the left.
Type of the indexer differs based on the type.
The return type differs based on the type.

| Type      | Indexing Type | Return Type | Returns                                                                                                              |
| --------- | ------------- | ----------- | -------------------------------------------------------------------------------------------------------------------- |
| `str`     | `int`         | `str`       | the len-1 string at that index, 0-indexed. Negative indexes supported. Error if the index is larger than the length. |
| `list[T]` | `int`         | `T`         | the element at that index, 0-indexed. Negative indexes supported. Error if the index is larger than the length.      |
| `map[K]V` | `K`           | `V`         | the value associated with the key. Error if the key doesn't exist in the map.                                        |

## Collection Types

`str`, `list`, and `map` are all collection types.
They can be iterated upon either through the manual `for ... { if ... then { break; }; };` construct, or the more ergonomic `for ... in collection { ... };`.

The `for ... in collection` construct expects at least one `var` declaration in the `...`.
The collection to be iterated over is passed as `collection`.
The first `var` found is used to store the collection elements on each iteration:

| Collection | Element Type |
| ---------- | ------------ |
| `str`      | `str`        |
| `list[T]`  | `T`          |
| `map[K]V`  | `K`          |