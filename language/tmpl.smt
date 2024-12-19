use io;

let name = io.read("What is your name? ");

let template = <>
    <h1>Hello {{ name }}!</h1>
</>;

io.writeln(template);
