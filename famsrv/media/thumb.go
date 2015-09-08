package media

import (
	_ "image/jpeg" // moo
	"os"
	"os/exec"
	"path/filepath"
)

var tmpPath = filepath.Join("tmp", "thumb")

// FindThumb creates a thumbnail for a media item with path relative to root path
func FindThumb(relpath string) (path string, err error) {
	err = validate(relpath)
	if err != nil {
		return
	}

	src := filepath.Join(rootPath, relpath)
	if err != nil {
		return
	}

	err = os.MkdirAll(tmpPath, 0700)
	if err != nil {
		return
	}

	dst := filepath.Join(tmpPath, relpath)

	/* skip creation if file already exists */
	/* FIXME: do date/crc check */
	if _, err = os.Stat(dst); err == nil {
		return dst, nil
	}

	var args = []string{
		"convert",
		"-define", "jpeg:preserve-settings",
		"-resize", "500x",
		src,
		dst,
	}

	var cmd *exec.Cmd
	gmPath, err := exec.LookPath("gm")
	if err != nil {
		return
	}
	cmd = exec.Command(gmPath, args...)
	err = cmd.Run()
	return dst, err
}

// OpenThumb returns a file descriptor to a thumbnail for media file with
// path relative to root path
func OpenThumb(relpath string) (f *os.File, err error) {
	path, err := FindThumb(relpath)
	if err != nil {
		return
	}

	return os.Open(path)
}

// OpenThumb returns a file descriptor to a thumbnail for this media item
func (item Item) OpenThumb() (*os.File, error) {
	return OpenThumb(item.Name)
}
