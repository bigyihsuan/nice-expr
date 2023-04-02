var x is int 10;

var y is int 0;
print("y=",y,"\n");
set y is x;
print("y=",y,"\n");
set y * y;  // y = 100
print("y=",y,"\n");
set y + 35; // y = 135
print("y=",y,"\n");
set y % 5;  // y = 0
print("y=",y,"\n");

var b is bool false;
print("b=",b,"\n");
set b is 10 % 2 = 0;
print("b=",b,"\n");
set b is not b;
print("b=",b,"\n");