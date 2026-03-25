# Go-http-server

A minimal HTTP/1.1 server built from scratch in Go — no frameworks, no external dependencies.

## About

This project started as a way to get a better grasp of network fundamentals and Go simultaneously. Rather than reaching for a framework, everything here is built on raw TCP connections and manual HTTP parsing, which forces a deeper understanding of what's actually happening under the hood.

## How it works

The server listens on a TCP socket and handles each connection manually:

- Parses raw HTTP/1.1 requests by hand (request line, headers, body)
- Routes requests based on the path
- Constructs well-formed HTTP responses with correct headers

## Running

```
go run .
```

The server starts on port `4221` by default.

## Status

Work in progress — more features and endpoints to be documented as the project grows.
