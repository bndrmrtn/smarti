# Smarti - More than a simple templating engine

Smarti is a templating engine that is designed to help backend developers deliver content to the frontend faster.
Smarti is cross-platform language. It is designed to be used anywhere with anything.
It helps you to make html templates with a simple syntax without learning each templating engine's syntax.

## Features

- Simple syntax
- Cross-platform
- Easy to use
- Fast
- Lightweight
- No dependencies

## Installation

```bash
go install github.com/smlgh/smarti@latest
```

## Usage

```bash
smarti run main.smt
# Or
smarti server # To start the server on port 8080 where you can pass data to the template language.
```

## Demo

### Simple code example

```smarti
use io;

let name = io.read("What is your name? ");
name = capitalize(name);

writef("Hello, %s!\n", name);
```

### Our Goal

```smarti
use request; // import the request module to get request data
use response as rw; // import the response module to send data to the client  (rw is an alias for response writer)

if request.method != "POST" {
  rw.status(405); // Set the status code to 405 and send the response
  return;
}

for let i = 0; i < api.data.length; i++ {
  let data = api.data[i];
  write(<>
    <h1>{{ data.title }}</h1>
    <p>{{ data.content }}</p>
  </>);
}

wr.status(200); // Set the status code to 200 and send the response
```

This code is our goal. We want to make a simple template language that can be used in any project.
We're working on it. We're trying to make it as good as possible.

## Error handling

Smarti does not have try-catch blocks.
If you want to handle errors, you can use a special macro called `@err`.
Err macro creates a temporary variable called `err` that contains the error message.

```smarti
use json;

let data = json.from('{"name": "John"}') @err;
if err != nil {
  writef("Error: %v\n", err);
  return;
}
```

You can use the `@err` macro anywhere in the code.
Errors can be omitted if you don't want to handle them.
In this case non-fatal operations can return non-initialized values.
This feature is still in development.
