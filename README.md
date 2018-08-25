# go-consistent

Source code analyzer that helps you to make your Go programs more consistent.

## Quick start / Installation

This installs `go-consistent` binary under your `$GOPATH/bin`:

```bash
go get github.com/Quasilyte/go-consistent
```

If `$GOPATH/bin` is under your system `$PATH`, `go-consistent` command should be available after that.
This should print the help message:

```bash
go-consistency --help
```

You can pass package names and separate Go filenames to the `go-consistent` tool:

```bash
go-consistency foo.go bar.go mypkg
```

Arguments like `...` are planned, but are not implemented yet.