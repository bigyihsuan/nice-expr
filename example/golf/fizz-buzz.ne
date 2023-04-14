// https://codegolf.stackexchange.com/questions/58615/1-2-fizz-4-buzz

for var i is int in range(1,101,1){printline(if i%15=0then{return"FizzBuzz";}else if i%3=0then{return"Fizz";}else if i%5=0then{return"Buzz";}else{return i;});};

// ungolfed

for var i is int in range(1,101,1) {
    printline(
        if i%15=0 then{
            return"FizzBuzz";
        } else if i%3=0 then {
            return"Fizz";
        } else if i%5=0 then {
            return"Buzz";
        } else {
            return i;
        }
    );
};