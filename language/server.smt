use request;

if request.method != "GET" {
    send("Bad request");
    setStatus(400);
    return;
}

send("Hello, world!");
