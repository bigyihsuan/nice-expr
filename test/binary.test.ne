// arithmetic operators
// numbers
println(1 + 1,); // 2
println(2 - 2,); // 0
println(2.2 * 3.0,); // 6.6
println(4 / 4,); // 1
println(4 % 3,); // 1

// strings
println("hello" + " " + "world",); // "hello world"
println("hello" - "l",); // "heo"
println("hello" - "e",); // "hllo"

// list
var l is list[int] [1,2,3,4,5,];
println(l + [6,],); // [1,2,3,4,5,6]
println(l - [3,],); // [1,2,4,5]

// logical operators
println(true and true,); // true
println(false or true,); // true

// comparison operators
println(1 = 1,); // true
println(2 > 1,); // true
println(1 >= 1,); // true
println(1 < 2,); // true
println(1 <= 1,); // true

// indexing operator
const s is str "abc123";
println(s_4,); // 2
println(l_4,); // 5
