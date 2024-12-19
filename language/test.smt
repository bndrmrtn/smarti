use io;
use strs;

let name = io.read("Waht is your name? ");
name = capitalize(name);


io.writef("Hello, %s!\n", name);
