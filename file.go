package rice

import (
	"os"
)

// File abstracts file methods so the user doesn't see the difference between VirtualFile and os.File
// This interface also implements the io.Reader, io.Seeker, io.Closer and http.File interfaces
// DEPRECATED, not being used anymore as Open(rice.File, error) != Open(http.File, error), even though both interfaces have equal methods
type File interface {
	Close() error
	Stat() (os.FileInfo, error)
	Readdir(count int) ([]os.FileInfo, error)
	Read([]byte) (int, error)
	Seek(offset int64, whence int) (int64, error)
}
