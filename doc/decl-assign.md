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
set x + " and me too!"; // string append:  x = "hello, world! and me too!"
set x - "aeiou"; // string set difference: x = "hll, wrld! nd m t!"

var y is int 5;
set y * 10; // y = 50
set y + 4;  // y = 54
set y % 9;  // y = 6
```

## Declarations, Assignments, and Blocks

By default, when you declare a variable,
it goes into a global "scope" that keeps track of the variables and constants that exist.
When the interpreter enters a block, it creates a new scope,
where identifiers defined outside of it are accessible.

You can declare variables with names that already exist in the outer context,
and (inside the inner scope) the name will be overwritten.
When doing assignments, a copy of the outer variable is made for the inner scope,
and then changed to the new value.

In the below example, in the first block, the variables `foo` and `bar` are

```cs
var quux is int 2;
var foo is int {return quux + 4;}; // 6
var bar is dec {return 0.5 + foo / quux;}; // 3.5
var baz is dec {
    var foo is dec 10.1; // is 10.1 in this block only
    var bar is dec 2.0; // is 2.0 in this block only
    return foo * bar;
}; // 20.2

{
    set foo is 1234; // can reassign outside variables, but they won't change the outside value
    println("inside anonymous block", foo);
}; // block not attached to anything

println("at the end", quux, foo, bar, baz);
```
