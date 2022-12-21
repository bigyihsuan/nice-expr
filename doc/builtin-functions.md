# Built-in Functions

## `print(any...)`, `println(any...)`

`print`, and `println` print their arguments to STDOUT.
For `print`, the elements are printed without spaces in between.
`println` on the other hand, puts a newline between each element,
as well as a trailing newline.

## `len(collection) int`

`len` takes in a `str`, `list`, or `map`,
and returns the number of elements in the collection.
For `str`, it returns the number of characters;
for `list` and `map` it returns the number of elements.

## `range(start int, end int, step int) list[int]`

`range` takes in 3 integers: the start of the range, the end of the range, and the step value.
It returns a list, containing every `step`-th integer in the interval `[start, end)`.
So if you call `range(1,11,2)`, it will return `[1,3,5,7,9]` (every odd number from 1 (inclusive) to 11 (exclusive)).

The step must be non-0.
If `start < end` and `step` is negative, an empty list will be returned.
