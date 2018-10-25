# go-consistent

Source code analyzer that helps you to make your Go programs more consistent.

## Quick start / Installation

This install `go-consistent` binary under your `$GOPATH/bin`:

```bash
go get github.com/Quasilyte/go-consistent
```

If `$GOPATH/bin` is under your system `$PATH`, `go-consistent` command should be available after that.<br>
This should print the help message:

```bash
go-consistent --help
```

You can pass package names and separate Go filenames to the `go-consistent` tool:

```bash
go-consistent foo.go bar.go mypkg
```

You can also use `std`, `./...` and other conventional targets that are normally
understood by Go tools.

* If you want to check consistency of a single file or package, just provide their name
* If you want to check the whole project, you should pass **all** its packages as an arguments

## Overview

To understand what `go-consistent` does, take a look at these 3 lines of code:

```go
lit1 := map[K]V{}
lit2 := map[K]V{}
m := make(map[K]V)
```

`lit1`, `lit2` and `m` are initialized to an empty, non-nil map.
The problem is that you have at least 2 ways to do it:

1. `lit1` and `lit2` use the first option, the map literal
2. `m` uses the second option, the `make` function

Neither of these are the "best", but on the package or project level, you might want to prefer
only one of them, for consistency reasons.

`go-consistent` tool detects that map literal used more frequently (2 vs 1) in the example above,
so it suggest you to replace `m` initialization expression to use map literal instead of `make` function.

There are many similar cases where you have 2 or more options of expressing the same thing in Go,
`go-consistent` tries to handle as much patterns as possible.

### Project traits

* Zero-configuration. Defaults should be good enough for most users.
  Other configuration is performed using command line arguments.
* Can handle any size of projects. This means that there should be no significant
  memory consumption growth with the increased number of packages being checked.
  There can be "fast, but memory-hungry" option that can work best for small-average projects,
  but it should be always possible to check huge projects on the developer machine.
