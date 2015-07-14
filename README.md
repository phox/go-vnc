# VNC Library for Go
go-vnc is a VNC library for Go, initially supporting VNC clients but
with the goal of eventually implementing a VNC server.

This library implements [RFC 6143][RFC6143].

## Project links
* Build Status:  [![Build Status][CIStatus]][CIProject]
* Documentation: [![GoDoc][GoDocStatus]][GoDoc]
* Views:         [![views][SGViews]][SGProject] [![views_24h][SGViews24h]][SGProject]
* Users:         [![library users][SGUsers]][SGProject] [![dependents][SGDependents]][SGProject]

## Setup
1. Download software and supporting packages.

    ```
    $ go get github.com/kward/go-vnc
    $ go get golang.org/x/net
    ```

## Usage
Sample code usage is available in the GoDoc.

- Connect and listen to server messages: <https://godoc.org/github.com/kward/go-vnc#example-Connect>

The source code is laid out such that the files match the document sections:

- [7.1] handshake.go
- [7.2] security.go
- [7.3] initialization.go
- [7.4] pixel_format.go
- [7.5] client.go
- [7.6] server.go
- [7.7] encodings.go

There are two additional files that provide everything else:

- vncclient.go -- code for instantiating a VNC client
- common.go -- common stuff not related to the RFB protocol


<!--- Links -->
[RFC6143]: http://tools.ietf.org/html/rfc6143

[CIProject]: https://travis-ci.org/kward/go-vnc
[CIStatus]: https://travis-ci.org/kward/go-vnc.png?branch=master

[GoDoc]: https://godoc.org/github.com/kward/go-vnc
[GoDocStatus]: https://godoc.org/github.com/kward/go-vnc?status.svg

[SGProject]: https://sourcegraph.com/github.com/kward/go-vnc
[SGDependents]: https://sourcegraph.com/api/repos/github.com/kward/go-vnc/.badges/dependents.svg
[SGUsers]: https://sourcegraph.com/api/repos/github.com/kward/go-vnc/.badges/library-users.svg
[SGViews]: https://sourcegraph.com/api/repos/github.com/kward/go-vnc/.counters/views.svg
[SGViews24h]: https://sourcegraph.com/api/repos/github.com/kward/go-vnc/.counters/views-24h.svg?no-count=1
