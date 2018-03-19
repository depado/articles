title: Copy files and directories in Go
description: Because there is no built-in recursive directory copy in Go
slug: copy-files-and-directories-in-go
author: Depado
date: 2016-09-29 14:00:00
tags:
    - dev
    - go

# Introduction

This article is intented to help you with the copy process of one or multiple files.
We'll start by creating a new package, and call it `copy`.
Our functions will be named respectively `File(src, dst string)` and
`Dir(src, dst string)` so that we can import our package and just call something like
`copy.File("file.txt", "file_copy.txt")`

# Copy a single file

```go
package copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

```


# Copy a directory recursively

```go
// Dir copies a whole directory recursively
func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
```
