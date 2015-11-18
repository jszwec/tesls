tesls [![GoDoc](https://godoc.org/github.com/jszwec/tesls?status.svg)](http://godoc.org/github.com/jszwec/tesls) [![Build Status](https://travis-ci.org/jszwec/tesls.svg)](https://travis-ci.org/jszwec/tesls)
==========

Lists tests in the given Go package without running them.

Installation
------------

    go get github.com/jszwec/tesls/cmd/tesls
    go install github.com/jszwec/tesls/cmd/tesls

Usage
-----
```
    tesls .
    > tesls     TestTests     /Users/jszwec/src/github.com/jszwec/tesls/tests_test.go
```

```
    tesls ./...
    > main      TestDirs      /Users/jszwec/src/github.com/jszwec/tesls/cmd/tesls/tesls_test.go
    > tesls     TestTests     /Users/jszwec/src/github.com/jszwec/tesls/tests_test.go
```

```
    tesls github.com/jszwec/tesls
    > tesls     TestTests     /Users/jszwec/src/github.com/jszwec/tesls/tests_test.go
```

```
    tesls github.com/jszwec/tesls/...
    > main      TestDirs      /Users/jszwec/src/github.com/jszwec/tesls/cmd/tesls/tesls_test.go
    > tesls     TestTests     /Users/jszwec/src/github.com/jszwec/tesls/tests_test.go
```

```
    tesls github.com/jszwec/tesls github.com/jszwec/tesls/cmd/tesls
    > main      TestDirs      /Users/jszwec/src/github.com/jszwec/tesls/cmd/tesls/tesls_test.go
    > tesls     TestTests     /Users/jszwec/src/github.com/jszwec/tesls/tests_test.go
```

```
    tesls -f='json' github.com/jszwec/tesls
    > [{"name":"TestTests","file":"/Users/jszwec/src/github.com/jszwec/tesls/tests_test.go","pkg":"tesls"}]
```

```
    tesls -f='Pkg: {{.Pkg}} | TestName: {{.Name}} | File: {{.File}}' github.com/jszwec/tesls
    > Pkg:     main      |     TestName:     TestDirs      |     File:     /Users/jszwec/src/github.com/jszwec/tesls/cmd/tesls/tesls_test.go
    > Pkg:     tesls     |     TestName:     TestTests     |     File:     /Users/jszwec/src/github.com/jszwec/tesls/tests_test.go
```

```
    tesls -tabs=false ./...
    > main TestDirs /Users/jszwec/src/github.com/jszwec/tesls/cmd/tesls/tesls_test.go
    > tesls TestTests /Users/jszwec/src/github.com/jszwec/tesls/tests_test.go
```
