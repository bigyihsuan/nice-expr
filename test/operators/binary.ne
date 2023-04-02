println("numbers");
println(1 + 1,); // 2
println(2 - 2,); // 0
println(2.2 * 3.0,); // 6.6
println(4 / 4,); // 1
println(4 % 3,); // 1
println("---");

println("strings");
println("hello" + " " + "world",); // "hello world"
println("hello" - "l",); // "heo"
println("hello" - "e",); // "hllo"
println("hello" - "el",); // "hlo"
println("---");

println("list");
var l is list[int] [1,2,3,4,5,];
println(l + [6,],); // [1,2,3,4,5,6]
println(l - [3,],); // [1,2,4,5]
println("---");

println("logical operators");
println(true and true,); // true
println(false or true,); // true
println("---");

println("comparison operators");
println(1 = 1,); // true
println(2 > 1,); // true
println(1.1 >= 1.0,); // true
println(0.3 < 2.2,); // true
println(1 <= 1,); // true
println(1 != 1,); // false
println(2 < 1,); // false
println(1.1 <= 1.0,); // false
println(0.3 > 2.2,); // false
println(1 > 1,); // false
println("---");

println("indexing operator");
const s is str "abc123";
println(s_4,); // 2
println(s_-1,); // 3
println("---");

println(l_4,); // 5
println(l_-1,); // 5
println("---");

const m is map[str]int <|"a":1,"b":2,|>;
println(m_"a",); // 1
