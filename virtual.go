package rice

import (
	"errors"
	"os"
	"syscall"
)

// Error indicating some function is not implemented yet (but available to satisfy an interface)
var ErrNotImplemented = errors.New("not implemented yet")

// virtualFile implements rice.File, it requires the box to be embedded (otherwise, the box should use os package to read files from disk)
type virtualFile struct {
	*EmbeddedFile
	offset int64 // read position on the virtual file
}

func newVirtualFile(ef *EmbeddedFile) *virtualFile {
	vf := &virtualFile{
		EmbeddedFile: ef,
		offset:       0,
	}
	return vf
}

func (vf *virtualFile) Close() error {
	vf.EmbeddedFile = nil
	vf.offset = 0
	return nil
}

func (vf *virtualFile) Stat() (os.FileInfo, error) {
	return vf.EmbeddedFile, nil
}

func (vf *virtualFile) Readdir(count int) ([]os.FileInfo, error) {
	//++ wont work for a file
	return nil, ErrNotImplemented
}

func (vf *virtualFile) Read(bts []byte) (int, error) {
	end := vf.offset + int64(len(bts))
	n := copy(bts, vf.Content[vf.offset:end])
	vf.offset += int64(n)
	return n, nil
}

func (vf *virtualFile) Seek(offset int64, whence int) (int64, error) {
	var e error

	//++ TODO: check if this is correct implementation for seek
	switch whence {
	case os.SEEK_SET:
		//++ check if new offset isn't out of bounds, set e when it is, then break out of switch
		vf.offset = offset
	case os.SEEK_CUR:
		//++ check if new offset isn't out of bounds, set e when it is, then break out of switch
		vf.offset += offset
	case os.SEEK_END:
		//++ check if new offset isn't out of bounds, set e when it is, then break out of switch
		vf.offset = vf.EmbeddedFile.Size() - offset
	}

	if e != nil {
		return 0, &os.PathError{
			Op:   "seek",
			Path: vf.Filename,
			Err:  e,
		}
	}

	return vf.offset, nil
}

// virtualDir implements rice.File, it requires the box to be embedded (otherwise, the box should use os package to read files from disk)
type virtualDir struct {
	*EmbeddedDir
}

func newvirtualDir(ed *EmbeddedDir) *virtualDir {
	vf := &virtualDir{
		EmbeddedDir: ed,
	}
	return vf
}

func (vd *virtualDir) Close() error {
	vd.EmbeddedDir = nil
	return nil
}

func (vd *virtualDir) Stat() (os.FileInfo, error) {
	return vd.EmbeddedDir, nil
}

func (vd *virtualDir) Readdir(count int) ([]os.FileInfo, error) {
	//++ read ChildDirs and ChildFiles from vd.EmbeddedDir
	return nil, ErrNotImplemented
}

func (vd *virtualDir) Read(bts []byte) (int, error) {
	// wont work for a dir (right?)
	return 0, errors.New("doesnt work for dir (TODO: proper error such as os's error)")
}

func (vd *virtualDir) Seek(offset int64, whence int) (int64, error) {
	return 0, &os.PathError{
		Op:   "seek",
		Path: vf.Filename,
		Err:  syscall.EISDIR,
	}
}
