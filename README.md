# Nexus
Nexus is a light weight HTTP web server written in Go

## Getting Started
Drawing inpiration from the net/http package in the standard library, nexus aims to be an educative-oriented attempt to understand socket programming with Go
Nexus uses the net package provided by the standard library to create sockets and establish connections with clients.

## Learning Outcomes
At the moment this project has alot of learning opportunities for anyone who's interested in tinkering with the low-level networking stuff.
Among them are:
* Understanding socket programming from the Berkeley sockets implementation
* Tinkering with sockets and understanding what they are
* Understanding under the hood implementation of net/http

## Capabilities
The server is able to:
* 1. Listen to and accept connections from clients
* 2. Read request payload from the client and respond with an HTTP 200 OK status for successful requests

## Tests
As it stands, the project has a single test, which tests the server's capability to accept connection requests from a client and response with a HTTP 200 OK response.

## Improvements
I am continuously working on this project to make it a fully functional web server, with capabilities similar to most of the modern commercial servers.
Among areas that need improvements are:
* Better request handling, i.e parsing JSON requests successfully for different HTTP methods
* Respond with appropriate HTTP methods to show request success or failure
* Serving static files, i.e static web pages and content
* CGI-like capability to execute external scripts and respond with dynamic content

## Contributing
Contributions are generously welcomed. To get started, for the repository and create a clone of the repo
