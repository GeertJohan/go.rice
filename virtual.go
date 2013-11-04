package rice

import (
	"errors"
	"os"
)

// Error indicating some function is not implemented yet (but available to satisfy an interface)
var ErrNotImplemented = errors.New("not implemented yet")

// virtualFile implements rice.File, it requires the box to be embedded (otherwise, the box should use os package to read files from disk)
type virtualFile struct {
	box          *Box          //++ ?? is this required?
	embeddedFile *EmbeddedFile // pointer to embedded file
	pos          int64         // read position on the virtual file
}

func newVirtualFile(ef *EmbeddedFile) *virtualFile {
	vf := &virtualFile{
		box:          nil, //++ ?? is this required?
		embeddedFile: ef,
		pos:          0,
	}
	return vf
}

func (vf *virtualFile) Close() error {
	vf.embeddedFile = nil
	vf.pos = 0
	return nil
}

func (vf *virtualFile) Stat() (os.FileInfo, error) {
	return vf.embeddedFile, nil
}

func (vf *virtualFile) Readdir(count int) ([]os.FileInfo, error) {
	//++ do stuff with vf.box.embeddedBox
	return nil, ErrNotImplemented
}

func (vf *virtualFile) Read([]byte) (int, error) {
	//++ copy bytes from vi.embeddedFile.Content into buf
	return 0, ErrNotImplemented
}

func (vf *virtualFile) Seek(offset int64, whence int) (int64, error) {
	//++ find how to implement this
	return 0, ErrNotImplemented
}

// type virtualFileInfo struct {
// 	vf *virtualFile
// }

// func (vfi *virtualFileInfo) Name() string {
// 	return vfi.vf.embeddedFile.Name
// }

// func (vfi *virtualFileInfo) Size() int64 {
// 	return int64(len(vfi.vf.embeddedFile.Content))
// }

// func (vfi *virtualFileInfo) Mode() os.FileMode {
// 	//++ static filemode for all Virtual files/dirs?
// 	return os.FileMode(0)
// }

// func (vfi *virtualFileInfo) ModTime() time.Time {
// 	return vfi.vf.embeddedFile.ModTime
// }

// func (vfi *virtualFileInfo) IsDir() bool {
// 	//++ get from embed huh..
// 	return false
// }

// func (vfi *virtualFileInfo) Sys() interface{} {
// 	return nil
// }
