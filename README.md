intern
======

[![Go Report Card](https://goreportcard.com/badge/github.com/spakin/intern)](https://goreportcard.com/report/github.com/spakin/intern) [![Build Status](https://travis-ci.com/spakin/intern.svg?branch=master)](https://travis-ci.com/spakin/intern) [![Go Reference](https://pkg.go.dev/badge/github.com/spakin/intern.svg)](https://pkg.go.dev/github.com/spakin/intern)

Description
-----------

[_String interning_](https://en.wikipedia.org/wiki/String_interning) traditionally means mapping multiple strings with equal contents to a single string in memory.  Doing so improves program speed by replacing slow string comparisons with fast pointer comparisons.

The `intern` package for [Go](https://golang.org/) takes a slightly different approach: It maps strings to _integers_.  This serves two purposes.  First, in Go, unlike C, a string is not simply a pointer to an array of characters so the pointer-comparison approach doesn't apply as naturally in Go as it does in C.  Second, the package provides a rather unique mechanism for preserving order.  That is, strings can be mapped to integers in such a way that if one string precedes another alphabetically, then the first string will map to a smaller integer than the second string.

More specifically, `intern` provides two symbol types: `Eq`, which supports only comparisons for equality and inequality (`==` and `!=`); and `LGE`, which supports less than, greater than, and equal to comparisons (`<`, `>`, `==`, `<=`, `>=`, and `!=`).  The former is faster to allocate, and allocation always succeeds.  The latter is slower to allocate and generally requires a pre-allocation step to avoid running out of integers that preserve the less than/greater than/equal to properties between all pairs of interned strings.

Installation
------------

Assuming you're using the [Go module system](https://go.dev/blog/using-go-modules), simply import the package into your program with
```Go
import "github.com/spakin/intern"
```
and run
```bash
go mod tidy
```
to download and install `intern`.

Documentation
-------------

Descriptions and examples of the `intern` API can be found online in the [pkg.go.dev documentation](https://pkg.go.dev/github.com/spakin/intern).

Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott+int@pakin.org*
