package rice

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	embeds = make(map[string]*EmbeddedBox) // maps box name to *EmbeddedBox
)

// Box abstracts a directory for resources/files.
// It can either load files from disk, or from embedded code (when `rice --embed` was ran).
type Box struct {
	name         string
	absolutePath string
	embed        *EmbeddedBox
}

// FindBox returns a Box instance for given name.
// When the given name is a relative path, it's base path will be the calling pkg/cmd's source root.
// When the given name is absolute, it's absolute. derp.
// Make sure the path doesn't contain any sensitive information as it might be placed into generated go source (embedded).
func FindBox(name string) (*Box, error) {
	b := &Box{
		name: name,
	}

	// find if box is embedded
	if embed := embeds[name]; embed != nil {
		b.embed = embed
		return b, nil
	}
	// box was not embedded

	// when given name is an absolute path, set it as absolute path.
	// otherwise calculate absolute path from caller source location
	if filepath.IsAbs(name) {
		// ++ think about abspath as name
		fmt.Println("probably shouldn't allow this.. rice.FindBox(..) should only take relative arguments, as the box is always inside the go package.. Might create FindBoxAbs(name, absPath) to work with box outside path.")
		os.Exit(-2)
		b.absolutePath = name
	} else {
		// resolve absolute directory path (when )
		err := b.resolveAbsolutePathFromCaller()
		if err != nil {
			return nil, err
		}
	}

	// check if absolutePath exists on filesystem
	info, err := os.Stat(b.absolutePath)
	if err != nil {
		return nil, err
	}
	// check if absolutePath is actually a directory
	if !info.IsDir() {
		return nil, errors.New("given name/path is not a directory")
	}

	// all done
	return b, nil
}

func (b *Box) resolveAbsolutePathFromCaller() error {
	_, callingGoFile, _, ok := runtime.Caller(2)
	if !ok {
		return errors.New("couldn't find caller on stack")
	}

	// resolve to proper path
	pkgDir := filepath.Dir(callingGoFile)
	b.absolutePath = filepath.Join(pkgDir, b.name)
	return nil
}

// AbsolutePath returns the absolute directory path for this box
// The path might not exist, it's only valid at the time and machine the calling command was compiled.
//++ TODO: should this method be exported? Whats the use? Should it be used? When it's being used, that defeats the purpose of this pkg..
func (b *Box) AbsolutePath() string {
	return b.absolutePath
}

// IsEmbedded indicates wether this box was embedded into the application
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
//   e.g.: http.Handle("/", http.FileServer(rice.Box("http-files")))
// If there is an error, it will be of type *os.PathError.
//++ TODO: don't return http.File, but return box.File if that qualifies for http.FileSystem
func (b *Box) Open(name string) (http.File, error) {
	if b.IsEmbedded() {
		name = strings.TrimLeft(name, "/")
		fmt.Printf("opening %s\n", name)
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
	file, err := os.Open(filepath.Join(b.absolutePath, name))
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
		cpy := make([]byte, 0, len(ef.Content))
		cpy = append(cpy, ef.Content...)
		// return copied bytes
		return cpy, nil
	}

	// open actual file from disk
	file, err := os.Open(filepath.Join(b.absolutePath, name))
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
		return ef.Content, nil
	}

	// open actual file from disk
	file, err := os.Open(filepath.Join(b.absolutePath, name))
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
