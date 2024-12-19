use response as rw;
use request as r;

let method = r.method();
let msg = format("Request method is: %s", method);

rw.header('Content-Type', 'text/html');

rw.write(<>
    <h1>Request method</h1>
</>);
