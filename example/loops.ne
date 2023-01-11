println("start of loop");
var sumis int 0;
for var i is int 0 {
    print("i=",i);
    println();
    print("s=",s);
    println();
    if i = 2 then {
        break;
    };
    set i += 1;
    set s += i;
};
println("after loop");
print("sum=",sum);
println();
