namespace server;
use request;
use response;
use httpsec;
use strs;

let method = request.method();

// return false if the method is not GET and write a response
if method != "GET" {
    // set header content type to text html
    response.header("Content-Type", "text/html");
    response.write(<>
        <h1>{{ method }} Method not allowed</h1>
    </>);
    return;
}

// get the name variable from the request query
let name = request.query("name");
if name == "" {
    name = "Unknown";
}

func greet(name) {
    // escape name to avoid XSS attacks
    name = httpsec.escapeHTML(name);
    return strs.concat("Hello, ", name, "!");
}

let message = greet(name);
message = <>
    <h1 style="color:#555;">{{ message }}</h1>
</>;

// set header content type to text html
response.header("Content-Type", "text/html");
response.write(message);

let style = <>
    <style>
        body {
            background: #222;
            font-family: Arial, Helvetica, sans-serif;
        }
    </style>
</>;

response.write(style);
