# Indexing

There are 2 indexing operators: `_` "get index" and `@` "set index".
`_` can be used anywhere to get a value at a specific index.

`@` is only used in `set` expressions to set at the index in the collection a value.
You may only use `is` in a `set ...@...` expression, but you can have any expression that returns a value after the `is`.

```swift
var l is list[int] range(1,6,1);
printline(l_3); // ok, prints `4`
printline(l@3); // not ok, only allowed in set
set l@3 is 10; // ok, sets the element at index 3 to be 10
printline(l_3); // should be 10
```

## Indexing Strings

Getting an index from a string returns another string,
but only a single character of it.
Setting an index to a string inserts the string into the target string:

```swift
var s is str "hello";
set s@2 is "world";
printline(s); // "heworldlo"
```

## Indexing Lists

Indexing a list will get the element at that index.
Setting an index will set the element at that index to the value.

## Indexing Maps

Indexing a map will get the element with the given key.
Setting an index will set the value of the given key to the given value.
If the key does not exist, it will be created.