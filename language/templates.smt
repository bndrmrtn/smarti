/*
    Smarti supports templates.
    Templates are spcecial variables that can hold html.
    Templates are opened with the <> operator and with the </> operator.
*/

// Create an html h1 template
let h1template = <>
    <h1>Hello, World!</h1>
</>;
// With this, you don't need to suffer with the string concatenation

// Print the template
writeln(h1template); // <h1>Hello, World!</h1>
