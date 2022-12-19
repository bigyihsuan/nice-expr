var x is int 9;
var z is int if x < 10 then {return x / 2;} else {return x * 2;};
println(z);

// standalone nested if-else
if x < 10 then {
  println("less than 10");
} else if x = 10 then {
  println("equal to 10");
} else {
  println("more than 10");
};

if z > 5 then {
    println("this should NOT print");
}