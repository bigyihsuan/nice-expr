printline("strings");
var s is str "hello";
printline(s); // hello
printline(s_1..-2); // el
set s@2 is "world";
printline(s); // heworldlo
printline(s_-2); // l
set s@-2 is "a";
printline(s);
printline(s_-2); // a

printline("---");
printline("lists");
var l is list[int] [2,4,6,8,10];
printline(l);
printline(l_2); // 6
printline(l_-1); // 10
printline(l_-2); // 8
printline(l_l_0); // 6
set l@2 is 0;
printline(l);
printline(l_2); // 0
printline(l_l_2); // 2

printline("---");
printline("maps");
var m is map[str]int <|"a":1, "b":2, "c":3|>;
printline(m);
printline(m_"c"); // 3
set m@"c" is 10 + m_"c";
printline(m_"c"); // 13
set m@"d" is 100;
printline(m_"d"); // 100
printline(m);

printline("---");
printline("slicing");
const r is list[int] range(0,11,1);
printline(r);
printline(r_0..len(r)); // entire list
printline(r_0..3); // [0,1,2]
printline(r_5..9); // [5,6,7,8]
printline(r_0..-1); // whole list excpet last element