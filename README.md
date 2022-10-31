# nice-expr

what if everything is an expression that return a value

even the loops

much of the lexing code is stolen from myself from [`bigyihsuan/nicer-syntax`](https://github.com/bigyihsuan/nicer-syntax/)

## Goals

* lots of keywords
* as few non-letter syllables as possible (?)
* everything is expression-based

## Features

* [x] Primitive literals (int, dec, str, bool)
* [x] Compound literals (list, map)
* [x] Type inference of compound literals
* [x] Arbitrarily-nested types
* [x] Declarations (var, const): nested, mixed, type checking
* [x] Function types: func(T...)V, funcs as return type
* [ ] Arithmetic operators (+, -, *, /, %)
* [ ] Logical operators (and, or, not)
* [ ] Assignments (set): check if var or const
* [ ] Arithmetic-Assignment operators (+=, -=, *=, /=, %=)
* [ ] Blocks: last un-semicolon-terminated return value
* [ ] Functions: return (early), last un-semicolon-terminated return value, recursion
* [ ] Conditionals: if-else, if-else if-else
* [ ] Loops: break, for, for-with-local-variable

## Examples?

```go
// assignments return the value of the variable
// assignments must have the type
var x int is 10;                  // returns 10
const n int is 10;                // also returns 10, but `n` cannot be have its value changed
const k int is (var y int is 10); // right-associative, parens are required
                                 // x = 10, y = 10
                                 // types must match
                                 // can mix const and var

// types
10          // int
1.23        // float
"a"         // string (no chars, just string)
[1,2,3,4,5] // list, homogenous
<|1:"a",2:"b",|> // map, key types must match, value types much match

// type zero-values
// int    -> 0
// float  -> 0.0
// string -> ""
// list   -> []
// map    -> <||>

// conditional expressions return the last expression in the option gone down
// all branches must return the same type
var z int is if x = 10 {
    x / 2
} else {
    x * 2
};

// loops
// infinite loop
for {};
// `break` stops the loop immediately, and returns `value`
for {break value};
for var i int is 0 { /* use i... */ }; // loop with loop variable
var x string is for { break 30 }; // x == 30

// example: getting and returning a sum
var l list[int] is [1,2,3,4,5,6,];
var sum is for var s int, var idx int {
    if idx = len(l) { // built-in function len() returns the length
        break s // return the sum when no more items
    }
    set s += l_idx // index lists using ints >=0 only
    set idx += 1
};

// all functions are anonymous until assigned to a variable
var greet func(string)string is func(name string)string{ "hello " + name };
// type of `greet` is func(string)string
var str string is greet("bob");

// recursion
var factorial func(int)int is func(n int)int {
    if n < 1 { return 1 } // early return
    if n = 2 { return 2 } // early return
    n * factorial(n-1) // recursive call
};
var num int is factorial(4);

// indexing
// use underscore to index string, list, map
// strings and list are 0-indexed
// maps take their key
// string -> string
// list_T -> T
// map[K]V -> V
"abcdefghij"_5 = "f"
[1,2,3,4,5,]_2 = 3
<|1:"a",2:"b",|>_1 = "a"
```

variables can only have letters
