# go-consistent

Source code analyzer that helps you to make your Go programs more consistent.

## Quick start / Installation

This installs `go-consistent` binary under your `$GOPATH/bin`:

```bash
go get github.com/Quasilyte/go-consistent
```

If `$GOPATH/bin` is under your system `$PATH`, `go-consistent` command should be available after that.<br>
This should print the help message:

```bash
go-consistenct --help
```

You can pass package names and separate Go filenames to the `go-consistent` tool:

```bash
go-consistenct foo.go bar.go mypkg
```

> Arguments like `...` are planned, but are not implemented yet.

`go-consistent` checks provided Go source code for consistency.

* If you want to check consistency of a single file or package, just provide their name
* If you want to check the whole project consistency, you should pass **all** its packages as an arguments
