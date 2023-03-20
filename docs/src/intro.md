# Introduction

wumpgo is a modular Discord Library for Golang.  Every component of the library can be used (almost) independently from the rest, and are mostly described by interfaces allowing you to provide your own implementations.  It was created out of frustration with existing libraries and a lack of ability to control the underlying behavior of the library, and hiding all of the inner workings behind abstractions.

wumpgo aims to expose as much or as little of the innerworkings as the developer desires.  For example, if you wish to have caching on your REST client, we provide two implementation, but both conform to an interface that you are welcome to implement yourself.

# Installation

The library is hosted on our own domain, so you can install it using

```bash
go install wumpgo.dev/wumpgo
```

# Documentation

This book aims to be the primary source of high level documentation, and a pure reference of the library itself is always availabe in the [Godoc](https://pkg.go.dev/wumpgo.dev/wumpgo).  This book is intended to be your guide to developing a bot, while the godoc will handle the nitty gritty of individual library components.