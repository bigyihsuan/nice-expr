var factorial is func(int)int func(var n is int) int {
    if n <= 2 then { return n; }; // early return
    return n * factorial(n-1);   // recursive call
};
var num is int factorial(4);
printline(num); // 24
