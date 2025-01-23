# Nexus

## Table of content
1. [About](#about)
2. [Getting started](#getting-started)
3. [Architecture](#architecture)
   1. [How nexus supports concurrent connection handling](#how-nexus-supports-concurrent-connection-handling)
4. [Feature checklist](#feature-checklist)
5. [Learning objectives](#learning-objectives)
6. [Reference](#reference)

# About
Nexus is a minimalistic light-weight static HTTP web server built with golang'.
It closely follows the HTTP protocol as described under the [RFC document](https://www.rfc-editor.org/rfc/rfc9110.html).

Eventhough inspired by Caddy and Nginx, Nexus doesn't aim to be as performant as Nginx or easily configurable as Caddy, rather it aim builds upon the fundamental
HTTP and TCP concepts, laying a good ground for someone who wants to learn how web servers work y supporting a small subset of features provided by a static web server.

# Getting Started
Currently, Nexus only supports static file serving with possibility of extension to support dynamic routing.
Also, I am aiming to provide configuration support such as Caddfile provided by Caddy server to enable users configure nexus as per their use case.

## Example usage
```shell
curl localhost:8080
```
> [!NOTE]
> Nexus runs on port `8080` by default; you might need to configure it and provide a suitable port to avoid conflict with curretly running programs.

```shell
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /index.html HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.81.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Server: nexus (Ubuntu)
< Content-Type: text/html
< Content-Length: 622
<
<!DOCTYPE html>
<html>
<head>
    <title>Welcome to nexus!</title>
    <style>
        body {
            width: 35em;
            margin: 0 auto;
            font-family: Tahoma, Verdana, Arial, sans-serif;
        }
    </style>
</head>
<body>
    <h1>Welcome to nexus!</h1>
    <p>If you see this page, the nexus web server is successfully installed and
        working. Further configuration is required.</p>
    <p><em>Thank you for using nexus.</e
```

All requests made to non-specific location are routed to `index.html` by default.

# Architecture

## How nexus supports concurrent connection handling
Nexus server requests similar to golang inbuild http server, by leveraging on goroutines. when the server starts, 3 it listens to connections on the address provided during startup. Three notifying channels(Done, Quit, and signalChan) are initialized. The Done channel is used to notify that the server has successfully started and it can begin accepting connections.

This is done because connection accepting is done by a goroutine, so we don't want any operation to proceed before the server is properly initialized which would cause the `connection refused` error whn clients try to connect.

The signalChan blocks the main goroutine from exiting until an os interrupt of shutdown signal is sent by the user. Once a shutown signal is received, the signalchan channel is closed sending a signal to the Quit channel to call the ShutDown method to close the server.

Each connection is handled separately in its own goroutine, similar to how golang http server handles connections.

# Feature checklist
Features currently supported:
 - [x] Handle client connections
 - [x] Parse client requests to obtain the HTTP method, request path and protocol
 - [x] Set appropriate HTTP headers in server responses
 - [x] Respond with appropriate HTTP method based on success or failure during request handling.
 - [x] Robust error handling and logging

## Learning Outcomes
At the moment this project has alot of learning opportunities for anyone who's interested in tinkering with the low-level networking stuff.
Among them are:
* Understanding socket programming from the Berkeley sockets implementation
* Tinkering with sockets and understanding what they are
* Understanding under the hood implementation of net/http

## Tests
As it stands, the project has a single test, which tests the server's capability to accept connection requests from a client and response with a HTTP 200 OK response.

## Improvements
I am continuously working on this project to make it a fully functional web server, with capabilities similar to most of the modern commercial servers.
Among areas that need improvements are:
* Better request handling, i.e parsing JSON requests successfully for different HTTP methods
* Respond with appropriate HTTP methods to show request success or failure
* CGI-like capability to execute external scripts and respond with dynamic content

## Contributing
Contributions are generously welcomed, whether bug fixes or feature requests. To get started, for the repository and create a clone of the repo
