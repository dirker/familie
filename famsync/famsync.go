package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func copyFile(srcFname string, dstFname string) (written int64, err error) {
	src, err := os.Open(srcFname)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.Create(dstFname)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	written, err = io.Copy(dst, src)
	return
}

func main() {
	destPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mediaItems := getIPhotoFiles("Website")

	for i, item := range mediaItems {
		fmt.Printf("\rCopying %d/%d", i+1, len(mediaItems))

		fname := path.Base(item.path)

		dst := path.Join(destPath, fname)
		copyFile(item.path, dst)

		if len(item.comment) > 0 {
			_ = ioutil.WriteFile(dst+"-comment.txt", []byte(item.comment), 0644)
		}
	}
	fmt.Println()
}
