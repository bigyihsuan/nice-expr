var sum is func(int)int func(var n is int)int{
    var s is int 0;
    var i is int 0;
    for {
        set s += i;
        set i += 1;
        print("s=",s," i=",i," n=",n," i>=n ->",i >= n);
        printline();
        if i >= n then {
            printline("breaking");
            break;
        };
    };
    printline("returning");
    return s;
};
// printline(sum(0));
// printline(sum(1));
// printline(sum(2));
// printline(sum(3));
// printline(sum(4));
// printline(sum(5));
printline(sum(6));