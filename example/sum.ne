var sum is func(int)int func(var n is int)int{
    var s is int 0;
    var i is int 0;
    for {
        set s += i;
        set i += 1;
        print("s=",s," i=",i," n=",n," i>=n ->",i >= n);
        println();
        if i >= n then {
            println("breaking");
            break;
        };
    };
    println("returning");
    return s;
};
// println(sum(0));
// println(sum(1));
// println(sum(2));
// println(sum(3));
// println(sum(4));
// println(sum(5));
println(sum(6));