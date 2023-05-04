// Implements a brainfuck interpreter, according to
// https://codegolf.stackexchange.com/questions/84/interpret-brainfuck
// - Cell size: 8bit unsigned. Overflow is undefined.
// - Array size: 30000 bytes (not circled)
// - Comments are everything not in +-.,[]<>
// - no EOF symbol

const code is str "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++";

var f is func(str)none func(var C is str)none{var T is list[int]repeat(0,30001);var p is int;var B is map[int]int;var b is list[int];for var i is int0{var c is str C_i;if c="["then{set b+[i];}else if c="]"then{set B@b_-1is i;set B@i is b_-1;set b is b_0..-1;};set i+1;if i>=len(C)then{break;};};for var i is int0{var c is str C_i;if c=">"then{set p+1;}else if c="<"then{set p-1;}else if c="+"then{set T@p is(T_p)+1;}else if c="-"then{set T@p is(T_p)-1;}else if c="."then{print(char(T_p));}else if c=","then{set T@p is ord(inputchar());}else if c="["and(T_p)=0then{set i is B_i;}else if c="]"and(T_p)!=0then{set i is B_i;};set i+1;if i>=len(C)then{break;};};};

f(code);

// ungolfed
var interpret is func(str)none func(var code is str)none {
    // there are neither 8-bit, nor unsigned, types in nice-expr.
    // a list of ints will have to do.
    var tape is list[int] repeat(0,30000);
    var ptr is int 0; // start on the start of the tape.

    var brackets is map[int]int <||>; // keys and values are indexes of brackets.

    var bracketStack is list[int] []; // stack that holds left bracket indexes
    for var i is int 0 { // first pass: match up brackets
        var c is str code_i;
        if c = "[" then {
            set bracketStack + [i]; // push a left bracket
        } else if c = "]" then {
            // pop a left bracket
            var left is int bracketStack_-1;
            set bracketStack is bracketStack_0..-1;
            var right is int i;
            set brackets@left is right;
            set brackets@right is left;
        };
        set i + 1;
        if i >= len(code) then { break; };
    };

    for var i is int 0 {
        var c is str code_i;
        if c = ">" then {
            set ptr + 1;
        } else if c = "<" then {
            set ptr - 1;
        } else if c = "+" then {
            set tape@ptr is (tape_ptr) + 1;
        } else if c = "-" then {
            set tape@ptr is (tape_ptr) - 1;
        } else if c = "." then {
            print(char(tape_ptr));
        } else if c = "," then {
            set tape@ptr is ord(inputchar());
        } else if c = "[" and (tape_ptr) = 0 then {
            set i is brackets_i; // jump to the matching right bracket
        } else if c = "]" and (tape_ptr) != 0 then {
            set i is brackets_i; // jump to matching left bracket
        };

        set i + 1;
        if i >= len(code) then { break; };
    };
};