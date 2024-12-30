# Golang Vanity URLs for Import Paths

This package implements a server that can respond to the golang import paths protocol and redirect `go tool` to the correct repository path so that we can host our go packages at `go.rotational.io` instead of `github.com/rotationalio`.

This server implements some Rotational specific tools. If you're interested in hosting your own vanity URLs for Go packages, I suggest reading [Making a Golang Vanity URL](https://medium.com/@JonNRb/making-a-golang-vanity-url-f56d8eec5f6c) by Jon Betti and using his [go.jonrb.io/vanity](https://pkg.go.dev/go.jonnrb.io/vanity) server as a starting place.
