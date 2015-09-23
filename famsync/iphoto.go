package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include <unistd.h>
#import <Cocoa/Cocoa.h>

void album_write_items_to_fd(const char *albumName, int fd)
{
  NSString *scriptSource = [NSString stringWithFormat:@"tell application \"iPhoto\" to properties of photos of album \"%s\"", albumName];
  NSAppleScript *script = [[NSAppleScript alloc] initWithSource:scriptSource];
  NSAppleEventDescriptor *d = [script executeAndReturnError:nil];

  unsigned int numItems = (unsigned int)[d numberOfItems];
  NSLog(@"photo library count: %u", numItems);

  for (int i = 1; i <= [d numberOfItems]; i++)
  {
    NSAppleEventDescriptor *photoDesc = [d descriptorAtIndex:i];

    if (!photoDesc) {
      NSLog(@"photoDesc empty :-(");
      continue;
    }

    NSString *imagePath    = [[photoDesc descriptorForKeyword:'ipth'] stringValue];
    NSString *originalPath = [[photoDesc descriptorForKeyword:'opth'] stringValue];
    NSString *comment      = [[photoDesc descriptorForKeyword:'pcom'] stringValue];

    NSString *csvLine = [NSString stringWithFormat:@"\"%@\",\"%@\",\"%@\"\n",
      imagePath,
      originalPath,
      comment
    ];

    const char *data = [csvLine UTF8String];
    write(fd, data, strlen(data));
  }
}

*/
import "C"
import (
	"encoding/csv"
	"os"
)

type mediaItem struct {
	path         string
	originalPath string
	comment      string
}

func getIPhotoFiles(albumName string) []mediaItem {
	r, w, _ := os.Pipe()
	defer r.Close()

	wfd := w.Fd()
	C.album_write_items_to_fd(C.CString(albumName), C.int(wfd))
	w.Close()

	csvReader := csv.NewReader(r)
	records, _ := csvReader.ReadAll()

	mediaItems := make([]mediaItem, len(records))

	for i, record := range records {
		mediaItems[i] = mediaItem{
			path:         record[0],
			originalPath: record[1],
			comment:      record[2],
		}
	}

	return mediaItems
}
