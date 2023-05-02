// Implements a brainfuck interpreter, according to
// https://codegolf.stackexchange.com/questions/84/interpret-brainfuck
// - Cell size: 8bit unsigned. Overflow is undefined.
// - Array size: 30000 bytes (not circled)
// - Comments are everything not in +-.,[]<>
// - no EOF symbol

// >	Move the pointer to the right
// <	Move the pointer to the left
// +	Increment the memory cell at the pointer
// -	Decrement the memory cell at the pointer
// .	Output the character signified by the cell at the pointer
// ,	Input a character and store it in the cell at the pointer
// [	Jump past the matching ] if the cell at the pointer is 0
// ]	Jump back to the matching [ if the cell at the pointer is nonzero

// there are neither 8-bit, nor unsigned, types in nice-expr.
// a list of ints will have to do.
var tape is list[int] repeat(0,30000);
var ptr is int 0; // start on the start of the tape.

var code is str inputall();
var brackets is map[int]int <||>; // keys and values are indexes of brackets.

var bracketStack is list[int] []; // stack that holds left bracket indexes
for var i is int 0 { // first pass: match up brackets
    var c is str code_i;
    // print(i, " ", c, "\n");
    if c = "[" then {
        set bracketStack + [i]; // push a left bracket
        // print("pushed left bracket at ", i, "\n");
    } else if c = "]" then {
        // printline("got right bracket");
        // pop a left bracket
        var left is int bracketStack_-1;
        set bracketStack is bracketStack_0..-1;
        var right is int i;
        set brackets@left is right;
        set brackets@right is left;
        // print("found bracket pair: ", left, " ", right, "\n");
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