int
float
string

list[int]
list[float]
list[string]
list[list[int]]
map[string]int
map[string]list[float]
map[string]list[list[float]]
map[list[list[float]]]int
map[list[list[float]]]list[int]
map[list[list[float]]]map[string]string

func(none,)none
func(string,)none
func(none,)float
func(float, string,)    map[float]string
func(map[int]string,    list[list[int]], )string
func(list[int], func(int,)int, )func(int,)int