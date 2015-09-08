package media

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var (
	rootPath string
)

// SetRoot sets the path of the repository of media files to use
func SetRoot(path string) {
	rootPath = path
}

// Item is a mediaItem
type Item struct {
	Name      string
	CreatedAt time.Time
	Comment   string
}

type byCreatedAt []Item

func (items byCreatedAt) Len() int {
	return len(items)
}

func (items byCreatedAt) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items byCreatedAt) Less(i, j int) bool {
	return items[i].CreatedAt.Before(items[j].CreatedAt)
}

// GetItems returns all items found in rootPath
func GetItems() (items []Item, err error) {
	dir, err := os.Open(rootPath)
	if err != nil {
		return
	}

	fi, err := dir.Stat()
	if !fi.Mode().IsDir() {
		panic("need directory")
	}

	names, err := dir.Readdirnames(0)
	if err != nil {
		return
	}

	for _, name := range names {
		item, err := NewItem(name)
		if err != nil {
			continue
		}

		items = append(items, item)
	}

	sort.Sort(sort.Reverse(byCreatedAt(items)))

	return items, nil
}

// NewItem creates a new media item
func NewItem(relpath string) (item Item, err error) {
	err = validate(relpath)
	if err != nil {
		return
	}

	path := filepath.Join(rootPath, relpath)

	f, err := os.Open(path)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if !fi.Mode().IsRegular() {
		return item, fmt.Errorf("not a regular file: %s", relpath)
	}

	x, err := exif.Decode(f)
	if err != nil {
		return
	}

	item.Name = relpath
	item.CreatedAt, err = x.DateTime()
	if err != nil {
		return
	}

	/* check for comment */
	comment, err := ioutil.ReadFile(path + "-comment.txt")
	if err != nil {
		return
	}

	item.Comment = string(comment)
	return
}

// Open returns a file handle to a media file with path relative to root path
func Open(relpath string) (*os.File, error) {
	err := validate(relpath)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(rootPath, relpath)
	return os.Open(path)
}

func validate(relpath string) error {
	relpath = filepath.Clean(relpath)
	if strings.Contains(relpath, "../") {
		return fmt.Errorf("invalid path")
	}

	/* for now, only support jpegs */
	ext := filepath.Ext(relpath)
	ext = strings.ToLower(ext)
	if ext != ".jpg" {
		return fmt.Errorf("filetype not supported: %s", ext)
	}

	return nil
}

// Open returns a file handle for this media item
func (item Item) Open() (*os.File, error) {
	return Open(item.Name)
}
