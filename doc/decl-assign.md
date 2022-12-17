# Declarations and Assignments

## Declarations

Declarations define a variable as existing within a program.
You can declare `var`iables or `const`ants.
All declarations must have a type.
Declarations by default are defined with a 0-value (see [Types](types.md).)
Declarations *must* come before assignment.
Declarations can be combined with assignments.

```cs
// var/const name is type (value)
const zeroInt is int; // = 0
var x is int 1;
const lst is list[str]; // = []
```

## Assignments

Assignments, preceded by `set`, allow for redefining variables.
They use `is` to redefine variables.
Other operators for `set` are available, mostly math operations.
You can only assign a value to variables (`var`), not constants.

```cs
var x is str "10";
set x is "hello, world!";
set x += " and me too!"; // string append:  x = "hello, world! and me too!"
set x -= "aeiou"; // string set difference: x = "hll, wrld! nd m t!"

var y is int 5;
set y *= 10; // y = 50
set y += 4;  // y = 54
set y %= 9;  // y = 6
```
