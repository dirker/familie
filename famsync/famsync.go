package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

const websitePath = "${HOME}/Sites/Familie"

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
	destPath := os.ExpandEnv(websitePath)
	if !exists(destPath) {
		panic("destination directory not found")
	}

	destPath = path.Join(destPath, "media")
	err := os.MkdirAll(destPath, 0770)
	if err != nil {
		panic("could not create destination directory")
	}

	mediaItems := getIPhotoFiles("Website")

	for _, item := range mediaItems {
		fname := path.Base(item.path)

		dst := path.Join(destPath, fname)
		copyFile(item.path, dst)

		if len(item.comment) > 0 {
			_ = ioutil.WriteFile(dst+"-comment.txt", []byte(item.comment), 0644)
		}
	}
}
