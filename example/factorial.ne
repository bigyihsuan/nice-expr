var factorial is func(int)int func(var n is int) int {
    if n < 1 then { return 1; }; // early return
    if n = 2 then { return 2; }; // early return
    return n * factorial(n-1);     // recursive call
};
var num is int factorial(4);
println(num);
