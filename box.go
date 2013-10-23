package rice

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	embeds = make(map[string]*EmbeddedBox)
)

// Box defines a abstracted set of resources/files
type Box struct {
	name  string
	dir   string
	embed *EmbeddedBox
}

// FindBox returns a Box instance
func FindBox(name string) (*Box, error) {
	b := &Box{
		name: name,
	}

	//++ see if box was embedded (RegisterEmbed(name, stuff..))

	// resolve absolute directory path
	err := b.resolveDirFromStack()
	if err != nil {
		return nil, err
	}

	//++ check if dir exists, error when not and not embedded

	return b, nil
}

func (b *Box) resolveDirFromStack() error {
	_, callingGoFile, _, ok := runtime.Caller(2)
	if !ok {
		return errors.New("couldn't find caller on stack")
	}

	// resolve to proper path
	pkgDir := filepath.Dir(callingGoFile)
	b.dir = filepath.Join(pkgDir, b.name)
	return nil
}

// Dir returns the absolute directory path for this box
// The path might not exist, it's only valid at the time and machine the calling command was compiled.
//++ TODO: should this method be exported?
func (b *Box) Dir() string {
	if b.dir == "" {
		b.resolveDirFromStack()
	}
	return b.dir
}

// IsEmbedded indicates wether this file was embedded into the application
func (b *Box) IsEmbedded() bool {
	return b.embed != nil
}

// Time returns how actual the box is.
// When the box is embedded, it's value is saved in the embedding code.
// When the box is live, this methods returns time.Now()
func (b *Box) Time() time.Time {
	if b.IsEmbedded() {
		return b.embed.Time
	}

	return time.Now()
}

// Open opens a File from the box
// Box implements http.FileSystem with this method.
// This allows the use of Box with a http.FileServer.
//   e.g.: http.Handle("/", http.FileServer(rice.Box("http-files").HTTPFileSystem()))
// If there is an error, it will be of type *os.PathError.
func (b *Box) Open(name string) (http.File, error) {
	if b.IsEmbedded() {
		ef := b.embed.Files[name]
		if ef == nil {
			return nil, &os.PathError{
				Op:   "open",
				Path: name,
				Err:  os.ErrNotExist,
			}
		}

		// box is embedded
		vf := newVirtualFile(ef)
		return vf, nil
	}

	// perform os open
	file, err := os.Open(filepath.Join(b.dir, name))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Bytes returns the content of the file with given name as []byte
func (b *Box) Bytes(name string) ([]byte, error) {
	// check if box is embedded
	if b.IsEmbedded() {
		// find file in embed
		ef := b.embed.Files[name]
		if ef == nil {
			return nil, os.ErrNotExist
		}
		// clone byteSlice
		cpy := make([]byte, len(ef.Content))
		copy(ef.Content, cpy)
		// return copied bytes
		return cpy, nil
	}

	// open actual file from disk
	file, err := os.Open(filepath.Join(b.dir, name))
	if err != nil {
		return nil, err
	}
	// read complete content
	bts, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// return result
	return bts, nil
}

// String returns the content of the file with given name as string
func (b *Box) String(name string) (string, error) {
	// check if box is embedded
	if b.IsEmbedded() {
		// find file in embed
		ef := b.embed.Files[name]
		if ef == nil {
			return "", os.ErrNotExist
		}
		// return as string
		return string(ef.Content), nil
	}

	// open actual file from disk
	file, err := os.Open(filepath.Join(b.dir, name))
	if err != nil {
		return "", err
	}
	// read complete content
	bts, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	// return result as string
	return string(bts), nil
}
