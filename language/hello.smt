namespace main; // main is used by default
use io;

func looping(inx) {
    if inx == 5 {
        return;
    }

    inx = inx + 1;
    looping(inx);
}

func main() {
    looping(0);
}
