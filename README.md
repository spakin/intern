intern
======

Description
-----------

[_String interning_](https://en.wikipedia.org/wiki/String_interning) traditionally means mapping multiple strings with equal contents to a single string in memory.  Doing so improves program speed by replacing slow string comparisons with fast pointer comparisons.

The `intern` package for [Go](https://golang.org/) takes a slightly different approach: It maps strings to _integers_.  This serves two purposes.  First, in Go, unlike C, a string is not simply a pointer to an array of characters so the pointer-comparison approach doesn't apply as naturally in Go as it does in C.  Second, the package provides a rather unique mechanism for preserving order.  That is, strings can be mapped to integers in such a way that if one string precedes another alphabetically, then the first string will map to a smaller integer than the second string.

More specifically, `intern` provides two symbol types: `Eq`, which supports only comparisons for equality and inequality (`==` and `!=`); and `LGE`, which supports less than, greater than, and equal to comparisons (`<`, `>`, `==`, `<=`, `>=`, and `!=`).  The former is faster to allocate, and allocation always succeeds.  The latter is slower to allocate and generally requires a pre-allocation step to avoid running out of integers that preserve the less than/greater than/equal to properties between all pairs of interned strings.

Installation
------------

Instead of manually downloading and installing `intern` from GitHub, the recommended approach is to ensure your `GOPATH` environment variable is set properly then issue a
```bash
go get github.com/spakin/intern
```
command.

Documentation
-------------

Descriptions and examples of the `intern` API can be found online in the [GoDoc documentation of package `intern`](https://godoc.org/github.com/spakin/intern).

Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott+int@pakin.org*
