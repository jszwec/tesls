tesls
==========

Lists tests in the given Go package without running them.

Installation
------------

    go get github.com/jszwec/tesls/cmd/tesls
    go install github.com/jszwec/tesls/cmd/tesls

Usage
-----

    tesls github.com/jszwec/tesls
    tesls -format='json' github.com/jszwec/tesls
    tesls -format='%P.%T %F' github.com/jszwec/tesls
