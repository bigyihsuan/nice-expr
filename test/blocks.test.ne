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