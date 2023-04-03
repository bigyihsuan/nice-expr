var quux is dec 2.0;
var foo is dec {return quux + 4.0;}; // 6
var bar is dec {return 0.5 + foo / quux;}; // 3.5
print("before baz block: foo=",foo,"\n");
var baz is dec {
    var foo is dec 10.1; // is 10.1 in this block only
    print("inside baz block: foo=",foo,"\n");
    var bar is dec 2.0; // is 2.0 in this block only
    return foo * bar;
}; // 20.2
print("after baz block: foo=",foo,"\n");

{
    set foo is 1234.5678; // if not redeclared, set can change outside variables
    printline("inside anonymous block");
    print("foo=",foo,"\n");
}; // block not attached to anything

printline("at the end");
print("quux=",quux,"\n");
print("foo=",foo,"\n");
print("bar=",bar,"\n");
print("baz=",baz,"\n");
