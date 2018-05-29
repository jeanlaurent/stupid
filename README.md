https://circleci.com/gh/jeanlaurent/stupid/tree/master

# Stupid

A stupid tool to write a portable `Makefile` for a `Go` project.

It provides simple commands to bridge the differences between the behaviors of the shells invoked from a typical `Makefile` by `GNU Make`.

All commands which manipulate files and directories should take linux-style paths as inputs, e.g. be made of forward slashes.

Available commands:
* [cp](#cp)
* [home](#home)
* [rm](#rm)

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
* copying is recursive
* intermediate directories are created
* permissions are replicated
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

Example:
```
stupid rm build/*.tar.gz electron/web
```
