# Stupid

[![CircleCI](https://circleci.com/gh/jeanlaurent/stupid/tree/master.svg?style=svg)](https://circleci.com/gh/jeanlaurent/stupid/tree/master)

[![codecov](https://codecov.io/gh/jeanlaurent/stupid/branch/master/graph/badge.svg)](https://codecov.io/gh/jeanlaurent/stupid)

A stupid tool to write a portable `Makefile` for a `Go` project.

It provides simple commands to bridge the differences between the behaviors of the shells invoked from a typical `Makefile` by `GNU Make`.

All commands which manipulate files and directories should take linux-style paths as inputs, e.g. be made of forward slashes.

Wildcard expansion for source files are performed for `*` and `?` meaning they can be used even if the underlying shell does not support them.

Available commands:
* [cp](#cp)
* [date](#date)
* [home](#home)
* [rm](#rm)
* [tar](#tar)
* [untar](#untar)

## Installation

```
go get -u github.com/jeanlaurent/stupid
```

Most systems have `GNU Make` available through their package manager.

For Windows the [GnuWin32 Make](http://gnuwin32.sourceforge.net/packages/make.htm) albeit a bit outdated has been known to work well.

## Commands

### cp
```
stupid cp SRCS DST
```
Copies files and directories listed in `SRCS` into `DST`, with the following behavior:
* existing files are overwritten
* permissions are replicated
* directories are copied recursively
* `SRCS` are globbed before processing
* `DST` is created if needed with all intermediate directories
* `DST` is a directory if any of the following is true:
  * `DST` already exists and is a directory
  * `DST` ends with a trailing slash
  * `SRCS` has more than one source
  * `SRCS` is a single source directory
* `DST` is a file if one of the following is true
  * `DST` already exists and is a file
  * `SRCS` is a single source file and `DST` does not end with a trailing slash

Example:
```
stupid cp web/readme.txt web/dist/* electron/web
```

### date
```
stupid date
```
Prints the current date with RFC3339.

### home
```
stupid home
```
Prints the home directory of the current user.

The exact format is platform dependent (e.g. the result can contain back or forward slashes or both) and may contain spaces, therefore it should be quoted when used.

Example:
```
stupid cp build/library.yaml "$(shell stupid home)/.tootool/"
```

### rm
```
stupid rm SRCS
```
Removes the files and directories listed in `SRCS`, with the following behavior:
* non existing sources are ignored
* directories are removed recursively
* `SRCS` are globbed before processing

Example:
```
stupid rm build/*.tar.gz electron/web
```

### tar
```
stupid tar SRCS DST
```
Creates a `DST` tar archive containing the files and directories listed in `SRCS`, with the additional behavior:
* directories are processed recursively
* intermediate directories for `DST` are created
* `SRCS` are globbed before processing
* if `DST` extension is `.tar.gz` or `.tgz` it also applies gzip compression
* if `DST` is `-` the archive is written to the standard output

Example:
```
stupid tar project.app readme.txt build/project-darwin.tar.gz
```

### untar
```
stupid untar SRC DST
```
Extracts from a `SRC` tar archive to the `DST` directory, with the following behavior:
* `DST` is created if needed with all intermediate directories
* `SRCS` are globbed before processing
* if `SRC` extension is `.tar.gz` or `.tgz` it also performs gzip decompression
* if `SRC` is `-` the archive is read from the standard input

Example:
```
stupid untar pony.tar.gz deps/github.com/ponies
```
