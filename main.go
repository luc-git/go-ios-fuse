/*
 * hellofs.go
 *
 * Copyright 2017-2022 Bill Zissimopoulos
 */
/*
 * This file is part of Cgofuse.
 *
 * It is licensed under the MIT license. The full license text can be found
 * in the License.txt file at the root of this project.
 */

package main

import (
	"fmt"
	"os"
	"path"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/afc"
	"github.com/winfsp/cgofuse/fuse"
)

var files []string

type Hellofs struct {
	fuse.FileSystemBase
}

var afcService *afc.Connection
var err error
var isdir map[string]bool

func (self *Hellofs) Opendir(path string) (int, uint64) {
	return 0, 0
}

func (self *Hellofs) Open(path string, flags int) (errc int, fh uint64) {
	return 0, 0
}

func (self *Hellofs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	if _, key := isdir[path]; key {
		stat.Mode = fuse.S_IFDIR | 0555
	} else {
		stat.Mode = fuse.S_IFREG | 0444
	}
	if path == "/" {
		stat.Mode = fuse.S_IFDIR | 0555
	}
	return 0
}

func (self *Hellofs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	return
}

func (self *Hellofs) Readdir(dirpath string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	files, err = afcService.ListFiles(path.Join("/", dirpath), "*")
	for _, f := range files {
		if f != ".." && f != "." {
			filetype, err := afcService.Stat(path.Join(dirpath, f))
			if err != nil {
				fmt.Printf(err.Error() + "\n")
				continue
			} else {
				fmt.Printf("SUCCESSSSSS!!")
			}
			if filetype.IsDir() {
				isdir[path.Join(dirpath, f)] = true
			} else {
				isdir[f] = false
			}
			fill(f, nil, 0)
		}
	}
	return 0
}

func main() {
	var device ios.DeviceEntry

	device, err = ios.GetDevice("")
	if err != nil {

	}
	/*if err.Error() == "error getting devicelist"{
		return
	}*/
	afcService, err = afc.New(device)
	/*if err.Error() == "fsync: connect afc service failed" {
		fmt.Printf(err.Error())
		return
	}*/
	if err != nil {

	}
	isdir = make(map[string]bool)
	hellofs := &Hellofs{}
	host := fuse.NewFileSystemHost(hellofs)
	host.Mount("", os.Args[1:])
}
