printline("numbers");
printline(1 + 1,); // 2
printline(2 - 2,); // 0
printline(2.2 * 3.0,); // 6.6
printline(4 / 4,); // 1
printline(4 % 3,); // 1
printline("---");

printline("strings");
printline("hello" + " " + "world",); // "hello world"
printline("hello" - "l",); // "heo"
printline("hello" - "e",); // "hllo"
printline("hello" - "el",); // "hlo"
printline("---");

printline("list");
var l is list[int] [1,2,3,4,5,];
printline(l + [6,],); // [1,2,3,4,5,6]
printline(l - [3,],); // [1,2,4,5]
printline("---");

printline("logical operators");
printline(true and true,); // true
printline(false or true,); // true
printline("---");

printline("comparison operators");
printline(1 = 1,); // true
printline(2 > 1,); // true
printline(1.1 >= 1.0,); // true
printline(0.3 < 2.2,); // true
printline(1 <= 1,); // true
printline(1 != 1,); // false
printline(2 < 1,); // false
printline(1.1 <= 1.0,); // false
printline(0.3 > 2.2,); // false
printline(1 > 1,); // false
printline("---");

printline("indexing operator");
const s is str "abc123";
printline(s_4,); // 2
printline(s_-1,); // 3
printline("---");

printline(l_4,); // 5
printline(l_-1,); // 5
printline("---");

const m is map[str]int <|"a":1,"b":2,|>;
printline(m_"a",); // 1
