println("start of loop");
var sumis int 0;
for var i is int 0 {
    print("i=",i,"\n");
    print("s=",s,"\n");
    if i = 2 then {
        break;
    };
    set i + 1;
    set s + i;
};
println("after loop");
print("sum=",sum,"\n");
