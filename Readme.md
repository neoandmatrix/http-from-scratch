# HTTP Server from Scratch

This project implements a basic HTTP server from scratch in Go, designed to handle HTTP/1.1 requests. It includes custom request parsing, response generation, and routing logic, all built without relying on Go's standard `net/http` package.

---

## Features

- **Custom HTTP Request Parsing**:
  - Parses HTTP request lines, headers, and body.
  - Supports chunked reading for efficient memory usage.

- **Custom HTTP Response Handling**:
  - Generates HTTP/1.1-compliant responses.
  - Includes default headers and customizable status codes.

- **Routing**:
  - Handles specific routes like `/yourproblem`, `/myproblem`, and `/httpbin/stream`.
  - Streams data from external APIs (e.g., `httpbin.org`).

- **Error Handling**:
  - Returns appropriate HTTP status codes for bad requests (`400`) and internal server errors (`500`).

---

