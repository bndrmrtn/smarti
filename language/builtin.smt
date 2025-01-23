namespace main;
use io;
use strs;

func string#length(s) {
    return strs.length(s);
}

let x = "Hello, ";

let l = x.length();

io.writeln(l);
