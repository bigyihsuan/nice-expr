println("standalone forloop");
for var i is int 0 {
    if i >= 10 then {
        break;
    };
    print("i=",i,"\n");
    set i + 1;
};

var l is list[int] range(1,8,1);
var sum is int for var s is int 0, var idx is int 0 {
    if idx >= len(l) then { // built-in function len() returns the length
        break s; // return the sum when no more items
    };
    set s + l_idx; // index lists using ints >=0 only
    set idx + 1;
};
println(sum); // 28