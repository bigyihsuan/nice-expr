var F is func(int)int func(var n is int)int{if n<3then{return n;}else{return n*F(n-1);};};

for var i is int in range(1,20,1) {
    printline(F(i));
};