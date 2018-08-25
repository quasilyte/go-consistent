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
* If you want to check the whole project, you should pass **all** its packages as an arguments

## Overview

To understand what `go-consistent` does, look at these 3 lines of code:

```go
a := map[int]string{}
b := map[int]string{}
c := make(map[int]string)
```

`a`, `b` and `c` are initialized to an empty, non-nil map.
The problem is that you have at least 2 ways to do it:

1. `a` and `b` use the first option, the map literal.
2. `c` uses the second option, the `make` function.

Neither of these are the "best", but on the package or project level, you might want to prefer
only one of them, for consistency reasons.

`go-consistent` tool detects that map literal used more frequently (2 vs 1) in the example above,
so it suggest you to replace `c` initialization expression to use map literal instead of `make` function.

There are many similar cases where you have 2 or more options of expressing the same thing in Go,
`go-consistent` tries to handle as much patterns as possible.
