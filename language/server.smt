namespace server;
use response;

let name = "World";

let html = <>
    <h1>Hello, {{ name }}!</h1>
</>;

response.write(html);
