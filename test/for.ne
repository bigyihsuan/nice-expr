var l is list[int] [1,2,3,4,5,6,7,8];
var sum is int for var s is int 0, var idx is int 0 {
    if idx >= len(l) then { // built-in function len() returns the length
        break s; // return the sum when no more items
    };
    set s += l_idx; // index lists using ints >=0 only
    set idx += 1;
};
println(sum); // 36