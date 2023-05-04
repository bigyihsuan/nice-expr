var x is int 9;
var z is int if x < 10 then {return x / 3 + 2;} else {return x * 2;}; // 5
printline(x, z); // 9 5

// standalone nested if-else
if x < 10 then {
  printline("less than 10");
} else if x = 10 then {
  printline("equal to 10");
} else {
  printline("more than 10");
};

if z > 10 then {
  printline("this should NOT print");
};

if x < 3 then {
  printline("this should NOT print");
} else {
  printline("more than 3");
};

// if with no return type returns None
printline(if z < 10 then {return;});

// if with decl
if var k is int z+x while k > 15 then {
  printline("this should NOT print", k);
} else {
  printline(k);
};