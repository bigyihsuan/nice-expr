# Built-in Functions

## `print(any...)`, `printline(any...)`

`print`, and `printline` print their arguments to STDOUT.
For `print`, the elements are printed without spaces in between.
`printline` on the other hand, puts a newline between each element,
as well as a trailing newline.

## `inputchar() string`, `inputline() string`, `inputall() string`

These are the 3 functions for getting input.
They take input from STDIN, return a string, and have varying behaviors:

- `inputchar` consumes and returns exactly 1 character.
- `inputline` consumes up to and including a newline, and returns it.
- `inputall` consumes all of STDIN and returns it.

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

## `repeat(ele any, count int) list[any]`

`repeat` returns a list containing `ele` repeated `count` times.

## `char(int)`, `ord(str)`

`char` converts a positive int (a Unicode code point) into a single-character string.

`ord` converts a single-character string into a positive int (its Unicode code point).