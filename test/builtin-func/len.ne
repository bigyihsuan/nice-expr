var l is list[int] [123,456,789];
var m is map[str]int <|"abc":1,"def":2|>;
var s is str "hello";

const a is list[int] [];
const b is map[int]int <||>;
const c is str "";

printline(len(l)); // 3
printline(len(m)); // 2
printline(len(s)); // 5

printline(len(a)); // 0
printline(len(b)); // 0
printline(len(c)); // 0
