var factorial is func(int)int func(n is int) int {
    if n < 1 { return 1 }; // early return
    if n = 2 { return 2 }; // early return
    n * factorial(n-1)     // recursive call
};
var num is int factorial(4);
println(num);
