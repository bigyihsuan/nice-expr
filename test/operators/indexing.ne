var s is str "hello";
set s@2 is "world";
printline(s); // "helworldlo"
printline(s_-2); // l
set s@-2 is "a";
printline(s);
printline(s_-2); // a

printline("---");
var l is list[int] [2,4,6,8,10];
printline(l);
printline(l_2); // 6
printline(l_-2); // 8
printline(l_l_0); // 6
set l@2 is 0;
printline(l_2); // 0
printline(l_l_2); // 2
printline(l);

printline("---");
var m is map[str]int <|"a":1, "b":2, "c":3|>;
printline(m);
printline(m_"c"); // 3
set m@"c" is 10 + m_"c";
printline(m_"c"); // 13
set m@"d" is 100;
printline(m_"d"); // 100
printline(m);
