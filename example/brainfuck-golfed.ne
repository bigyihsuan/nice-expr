// Implements a brainfuck interpreter, according to
// https://codegolf.stackexchange.com/questions/84/interpret-brainfuck
// - Cell size: 8bit unsigned. Overflow is undefined.
// - Array size: 30000 bytes (not circled)
// - Comments are everything not in +-.,[]<>
// - no EOF symbol

const code is str "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++";

var f is func(str)none func(var C is str)none{var T is list[int]repeat(0,30001);var p is int;var B is map[int]int;var b is list[int];for var i is int0{var c is str C_i;if c="["then{set b+[i];}else if c="]"then{set B@b_-1is i;set B@i is b_-1;set b is b_0..-1;};set i+1;if i>=len(C)then{break;};};for var i is int0{var c is str C_i;if c=">"then{set p+1;}else if c="<"then{set p-1;}else if c="+"then{set T@p is(T_p)+1;}else if c="-"then{set T@p is(T_p)-1;}else if c="."then{print(char(T_p));}else if c=","then{set T@p is ord(inputchar());}else if c="["and(T_p)=0then{set i is B_i;}else if c="]"and(T_p)!=0then{set i is B_i;};set i+1;if i>=len(C)then{break;};};};

f(code);