# goutils
A collection of useful go utilities and functions.

[![Build Status](https://travis-ci.org/Jatgam/goutils.svg?branch=master)](https://travis-ci.org/Jatgam/goutils)

### Version
Works with semantically versioned modules, according to [semver](https://semver.org/). The version is pulled from the git tags of the repository. Example: v0.1.0

The first step is to make sure the version is passed into the build:
```shell
GIT_VERSION=$(shell git describe --always --dirty)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
LDFLAGS=-ldflags "-X github.com/jatgam/goutils/version.gitVersion=${GIT_VERSION} -X github.com/jatgam/goutils/version.gitBranch=${GIT_BRANCH}"
go build ${LDFLAGS} ...
```

main usage:
```go
package main
import (
    "fmt"
    "os"
    "github.com/jatgam/goutils/version"
)
var (
    appVersion = &version.Info{}
)
func init() {
    appVersion.ParseVersion()
}
func main() {
    fmt.Println(os.Args[0], "version:", appVersion.GetVersionString())
}
```

### Queue and Double Linked List
Meant to be thread safe and support go concurrency.
