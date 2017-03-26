# Simple Guetzli

A Guetzli compression helper

`go get github.com/romainmenke/simple-guetzli`

This requires : [Guetzli](https://github.com/google/guetzli)

---

### Options

```
Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -q, --quality=95     Quality in units equivalent to libjpeg quality
      --verbose        Verbose mode
  -f, --force          Force recompression
      --force-quality  Force recompression if quality changed
  -t, --threads=3      Max concurrent threads
  -v, --version        Guetzli Version

Args:
  [<source>]  Source directory
  [<output>]  Output directory
  [<log>]     Log directory, the log is used to prevent duplicate compressions
```

---

### Why?

Guetzli is cpu intensive and waiting for builds / compiles /... is something we all like to avoid.
This little tool keeps a log of compressed files and skips those that have already been done.

Now you can safely watch a folder with images and compress only that what needs to be done.


In short :

- compresses an entire folder at once
- executes compressions in parallel
- tracks what has been compressed

Useful :

- cancel at any time, finished compressions will not need to be redone.
- manage max threads.

Enjoy!
