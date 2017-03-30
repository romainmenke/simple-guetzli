# Simple Guetzli

A Guetzli compression helper

`go get github.com/romainmenke/simple-guetzli`

or download from the latest [release](https://github.com/romainmenke/simple-guetzli/releases).

This requires : [Guetzli](https://github.com/google/guetzli)

---

### Why?

Guetzli is cpu intensive and waiting for builds / compiles /... is something we all like to avoid.
This little tool keeps a log of compressed files and skips those that have already been done.

Now you can safely watch a folder with images and compress only that what needs to be done.


In short :

- compresses an entire folder at once
- executes compressions in parallel
- tracks what has been compressed
- determines file changes and recompresses

Useful :

- cancel at any time, finished compressions will not need to be redone.
- manage max threads.
- max memory is divided by number of Guetzli instances (so 1000mb ram with 4 threads will give 250mb ram for each)

---

### Options

```
Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -q, --quality=95     Visual quality to aim for, expressed as a JPEG quality value. Default value is 95.
      --verbose        Print a verbose trace of all attempts to standard output.
  -m, --memlimit=6000  Memory limit in MB. Guetzli will fail if unable to stay under the limit. Default limit is 6000
      --nomemlimit     Do not limit memory usage.
  -f, --force          Force recompression
      --force-quality  Force recompression if quality changed
  -t, --threads=3      Max concurrent threads. Default limit is the number of threads for the cpu minus 1
  -v, --version        Guetzli Version

Args:
  [<source>]  Source directory
  [<output>]  Output directory
  [<log>]     Log directory, the log is used to prevent duplicate compressions
```

---
