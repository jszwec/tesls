tesls [![GoDoc](https://godoc.org/github.com/jszwec/tesls?status.svg)](http://godoc.org/github.com/jszwec/tesls) [![Build Status](https://travis-ci.org/jszwec/tesls.svg)](https://travis-ci.org/jszwec/tesls)
==========

Lists tests in the given Go package without running them.

Installation
------------

    go get github.com/jszwec/tesls/cmd/tesls
    go install github.com/jszwec/tesls/cmd/tesls

Usage
-----

    tesls .
    tesls ./...
    tesls github.com/jszwec/tesls
    tesls github.com/jszwec/tesls/...
    tesls github.com/jszwec/tesls github.com/jszwec/tesls/cmd/tesls
    tesls -f='json' github.com/jszwec/tesls
    tesls -f='{{.Pkg}}.{{.Name}} {{.File}}' github.com/jszwec/tesls
