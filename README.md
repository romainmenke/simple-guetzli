# Simple Guetzli

A Guetzli compression helper

`go get github.com/romainmenke/simple-guetzli`

This requires : [Guetzli](https://github.com/google/guetzli)

---

### Options

```
Flags:
      --help        Show context-sensitive help (also try --help-long and --help-man).
  -q, --quality=84  Quality in units equivalent to libjpeg quality
  -v, --verbose     Verbose mode

Args:
  [<source>]  Source directory
  [<output>]  Output directory
  [<log>]     Log directory, the log is used to prevent duplicate compressions
```

---

### Why?

Guetzli is cpu intensive and waiting for builds / compiles /... is something we all like to avoid.
This little tool keeps a log of compressed files and skips those that have already been done. Now you can safely add it to a watcher.

Enjoy!