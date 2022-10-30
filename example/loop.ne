var l list[int] is [1,2,3,4,5,6,];
var sum is for var s int, var idx int {
    if idx = len(l) { // built-in function len() returns the length
        break s // return the sum when no more items
    };
    set s += l_idx; // index lists using ints >=0 only
    set idx += 1;
};
println(sum);