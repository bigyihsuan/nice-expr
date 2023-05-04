const insertionSort is func(list[int])list[int] func(var L is list[int])list[int] {
    for var i is int 1 {
        for var j is int i {
            var temp is int L_(j-1);
            set L@j-1 is L_j;
            set L@j is temp; 
            set j - 1;
            if j <= 0 or (L_(j-1)) <= (L_j) then {break;};
        };
        set i + 1;
        if i >= len(L) then {break;};
    };
    return L;
};

var input is list[int] [99,-2,53,4,67,55,23,43,88,-22,36,45];

printline(input);
printline(insertionSort(input));