package rice

import (
	"fmt"
	"os"
	"time"
)

// EmbeddedBox defines an embedded box
type EmbeddedBox struct {
	Name    string                   // box name
	Time    time.Time                // embed time
	Files   map[string]*EmbeddedFile // embedded files by full path
	Dirs    map[string]*EmbeddedDir  // embedded dirs by full path
	RootDir *EmbeddedDir             // root directory for this box
}

// EmbeddedSingle defines an embedded single
type EmbeddedSingle struct {
	Name string        // single name
	Time time.Time     // embed time
	File *EmbeddedFile // embedded file
}

// EmbeddedDir holds the layout (directory structure) for a rice box
// EmbeddedDir implements os.FileInfo
type EmbeddedDir struct {
	Filename   string
	DirModTime time.Time
	ChildDirs  []*EmbeddedDir  // both implement os.FileInfo, combined they can be returned by virtualDir.Readdir()
	ChildFiles []*EmbeddedFile // both implement os.FileInfo, combined they can be returned by virtualDir.Readdir()
}

// Name returns the base name of the directory
// (implementing os.FileInfo)
func (ed *EmbeddedDir) Name() string {
	return ed.Filename
}

// Size always returns 0
// (implementing os.FileInfo)
func (ed *EmbeddedDir) Size() int64 {
	return 0
}

// Mode returns the file mode bits
// (implementing os.FileInfo)
func (ed *EmbeddedDir) Mode() os.FileMode {
	return os.FileMode(0555 | os.ModeDir) // dr-xr-xr-x
}

// ModTime returns the modification time
// (implementing os.FileInfo)
func (ed *EmbeddedDir) ModTime() time.Time {
	return ed.DirModTime
}

// IsDir returns the abbreviation for Mode().IsDir() (always true)
// (implementing os.FileInfo)
func (ed *EmbeddedDir) IsDir() bool {
	return true
}

// Sys returns the underlying data source (always nil)
// (implementing os.FileInfo)
func (ed *EmbeddedDir) Sys() interface{} {
	return nil
}

// EmbeddedFile defines an embedded file
// EmbeddedFile implements os.FileInfo
type EmbeddedFile struct {
	Filename    string // filename
	FileModTime time.Time
	Content     string
}

// Name returns the base name of the file
// (implementing os.FileInfo)
func (ef *EmbeddedFile) Name() string {
	return ef.Filename
}

// Size returns the length in bytes for regular files; system-dependent for others
// (implementing os.FileInfo)
func (ef *EmbeddedFile) Size() int64 {
	return int64(len(ef.Content))
}

// Mode returns the file mode bits
// (implementing os.FileInfo)
func (ef *EmbeddedFile) Mode() os.FileMode {
	return os.FileMode(0555) // r-xr-xr-x
}

// ModTime returns the modification time
// (implementing os.FileInfo)
func (ef *EmbeddedFile) ModTime() time.Time {
	return ef.FileModTime
}

// IsDir returns the abbreviation for Mode().IsDir() (always false)
// (implementing os.FileInfo)
func (ef *EmbeddedFile) IsDir() bool {
	return false
}

// Sys returns the underlying data source (always nil)
// (implementing os.FileInfo)
func (ef *EmbeddedFile) Sys() interface{} {
	return nil
}

// RegisterEmbeddedBox registers an EmbeddedBox
func RegisterEmbeddedBox(name string, box *EmbeddedBox) {
	if _, exists := embeds[name]; exists {
		panic(fmt.Sprintf("EmbeddedBox with name `%s` exists already", name))
	}
	embeds[name] = box
}
