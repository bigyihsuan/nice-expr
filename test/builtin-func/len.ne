var l is list[int] [123,456,789];
var m is map[str]int <|"abc":1,"def":2|>;
var s is str "hello";

const a is list[int] [];
const b is map[int]int <||>;
const c is str "";

println(len(l));
println(len(m));
println(len(s));

println(len(a));
println(len(b));
println(len(c));
