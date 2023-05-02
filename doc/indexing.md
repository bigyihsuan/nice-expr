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

## Slicing

You can use the double-dot operator `col_start..end` to get a slice of elements, on get-indexes only.
`start` is the inclusive starting index, and `end` is the exclusive ending index.
Negative indexes are permitted.
You cannot slice on a `map`.
If the number of desired elements `end-start` is larger than the number of elements in the collection, it will simply return everything after `start`.
If `start` is larger than `end`, it will return an empty string/list.

Slicing a string yields a string.
Slicing a list yields a list.

```swift
var l is list[int] range(1,6,1);
printline(l_0..-1); // [1,2,3,4]
printline(l_-3..-1); // [1,2,3,4]
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