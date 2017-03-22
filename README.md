# Simple Guetzli

A Guetzli compression helper

`go get github.com/romainmenke/simple-guetzli`

This requires : [Guetzli](https://github.com/google/guetzli)

---

### Options

- `-h`            : help
- `-source`       : source directory
- `-out`          : output directory
- `-log`          : log directory
- `-level`        : compression level
- `trailing args` : exclusion -> simple `must not contain` logic

---

### Why?

Guetzli is cpu intensive and waiting for builds / compiles /... is something we all like to avoid.
This little tool keeps a log of compressed files and skips those that have already been done. Now you can safely add it to a watcher.

In short :

- compresses an entire folder at once
- executes compressions in parallel
- tracks what has been compressed

Enjoy!
