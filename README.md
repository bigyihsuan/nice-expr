# nice-expr

what if everything is an expression that return a value

even the loops

~~much of the lexing code is stolen from myself from [`bigyihsuan/nicer-syntax`](https://github.com/bigyihsuan/nicer-syntax/)~~

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
  * [x] Zero Values when no value given
* [x] Function types: func(T,T,...)V, funcs as return type
* [x] Using variables
* [x] Operators
  * [x] Unary Operators
    * [x] `not`, `-`
  * [x] Binary Operators
    * [x] Arithmetic operators (+, -, *, /, %)
    * [x] Comparison operators (=, <, <=, >, >=)
    * [x] Logical operators (`and`, `or`)
    * [x] Indexing operator `_`
* [x] Assignments (set): check if var or const
* [x] Arithmetic/Comparison operators in assignments (`set i + 5; set b or true;`)
* [x] Built-in functions
  * [x] print()/printline()
  * [x] len(): lists, maps, strings
  * [x] range(): python-like 3-part range
  * [x] inputchar()/inputline()/inputall(): get 1 char/1 line/all chars of input as a string
* [x] Blocks
  * [x] return value with `return` keyword
* [/] Functions
  * [x] return value with `return` keyword
  * [x] recursion
  * [x] functions act as closures (inner functions can use outer functions' variables)
* [x] Conditionals:
  * [x] if-else
  * [x] if-else if-else
  * [x] return value with `return` keyword
* [x] Loops
  * [x] `break`
  * [x] for
  * [x] for-with-local-variables
* [ ] Type conversions `type(val)`

## Examples?

```cs
// assignments return the value of the variable
// assignments must have the type
var x is int 10;                  // returns 10
const n is int 10;                // also returns 10, but `n` cannot be have its value changed
const k is int (var y is int 10); // right-associative, parens are required
// x = 10, y = 10
// types must match
// can mix const and var

// types
10          // int
1.23        // dec
"a"         // str (no chars, just str)
true        // bool
[1,2,3,4,5] // list, homogenous
<|1:"a",2:"b"|> // map, key types must match, value types much match

// type zero-values
// int  -> 0
// dec  -> 0.0
// str  -> ""
// bool -> false
// list -> []
// map  -> <||>

// blocks denote a new scope, variables and constants are local to that block
// blocks can return a value if it has the `return` keyword
// execution of the block stops on the first `return` and continues on the outside of the block
var v is dec {return (3 + 4 + 5) / 2};

// conditional expressions return the last expression in the option gone down
// all branches must return the same type
var z is int if x = 10 then {return x / 2;} else {return x * 2;};

// standalone nested if-else
if x < 10 then {
  printline("less than 10");
} else if x = 10 then {
  printline("equal to 10");
} else {
  printline("more than 10");
};

// loops
// infinite loop
for {};
// `break` stops the loop immediately, and returns `value`
for {break value;};
for var i is int 0 { /* use i... */ }; // loop with loop variable
var x is int for { break 30; }; // x == 30

// example: getting and returning a sum
var l is list[int] [1,2,3,4,5,6,];
var sum is int for var s is int 0, var idx is int 0 {
    if idx = len(l) then { // built-in function len() returns the length
        break s; // return the sum when no more items
    };
    set s += l_idx; // index lists using ints >=0 only
    set idx += 1;
};

// all functions are anonymous until assigned to a variable
// can have const and var arguments
// they behave like regular variables: const can't be assigned to, var can
var greet is func(str)str func(const start is str, var name is str)str{
  set name = " Mr. " + name;
  return start + name;
};
// type of `greet` is func(str)str
var greeting is str greet("Hello", "Bob");
printline(greeting); // Hello Mr. Bob

// recursion
var factorial func(int)int is func(n int)int {
    if n < 1 { return 1; }; // early return
    if n = 2 { return 2; }; // early return
    return n * factorial(n-1); // recursive call
};
var num is int factorial(4);

// indexing
// use underscore to index str, list, map
// strs and list are 0-indexed
// maps take their key
// str -> str
// list_T -> T
// map[K]V -> V
"abcdefghij"_5 = "f"
[1,2,3,4,5,]_2 = 3
<|1:"a",2:"b",|>_1 = "a"
```

variables can only have letters
