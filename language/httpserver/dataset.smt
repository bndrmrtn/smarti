namespace main;
use io;

func loop(n) {
    if n == 0 {
        return;
    }
    io.writeln("n = ", n);
    loop(n - 1);
}

loop(5);
