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
go get github.com/smlgh/smarti@latest
```

## Usage

```bash
smarti run main.smt
# Or
smarti server # To start the server on port 8080 where you can pass data to the template language.
```

## Demo

```smarti
use api; // Import the api module to use data from the request;

api.statusCode(200); // Set the status code to 200;

for let i = 0; i < api.data.length; i++ {
  let data = api.data[i];
  write(<>
    <h1>{{ data.title }}</h1>
    <p>{{ data.content }}</p>
  </>);
}
```

This code is our goal. We want to make a simple template language that can be used in any project.
We're working on it. We're trying to make it as good as possible.
