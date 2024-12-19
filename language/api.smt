use request;
use responseWriter as w;
use json;

if request.method != "POST" {
    w.status(405);
    w.write("Method not allowed");
    return;
}

data = json.from(request.body) @err; // @err is a macro that checks if the previous expression has an error
if err != nil {
    w.status(400);
    w.write("Invalid JSON");
    return;
}
