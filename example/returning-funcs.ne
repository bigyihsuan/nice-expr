// add is a generator function
const add is func(int) func(int)int func(const x is int) func(int)int {
    return func(const n is int)int {
        return n + x;
    };
};

const addFive is func(int)int add(5);

printline(addFive(10)); // 15