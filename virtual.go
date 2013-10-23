package rice

import (
	"os"
	"time"
)

// virtualFile implements rice.File, it requires the box to be embedded (otherwise, the box should use os package to read files from disk)
type virtualFile struct {
	box          *Box
	embeddedFile *EmbeddedFile
	pos          int64 //++ position on the virtual file
}

func newVirtualFile(ef *EmbeddedFile) *virtualFile {
	//++ could this method return an error?
	//++ 	we expect EmbeddedFile to be existing, and this function is not exported..

	//++ setup vf, return vf
	return nil
}

func (vf *virtualFile) Close() error {
	//++
	return nil
}

func (vf *virtualFile) Stat() (os.FileInfo, error) {
	//++
	return nil, nil
}

func (vf *virtualFile) Readdir(count int) ([]os.FileInfo, error) {
	//++ do stuff with vf.box.embeddedBox
	return nil, nil
}

func (vf *virtualFile) Read([]byte) (int, error) {
	//++ copy bytes from vi.embeddedFile.Content into buf
	return 0, nil
}

func (vf *virtualFile) Seek(offset int64, whence int) (int64, error) {
	//++ find how to implement this
	return 0, nil
}

type virtualFileInfo struct {
	vf *virtualFile
}

func (vfi *virtualFileInfo) Name() string {
	return vfi.vf.embeddedFile.Name
}

func (vfi *virtualFileInfo) Size() int64 {
	return int64(len(vfi.vf.embeddedFile.Content))
}

func (vfi *virtualFileInfo) Mode() os.FileMode {
	//++ static filemode for all Virtual files/dirs?
	return os.FileMode(0)
}

func (vfi *virtualFileInfo) ModTime() time.Time {
	return vfi.vf.embeddedFile.ModTime
}

func (vfi *virtualFileInfo) IsDir() bool {
	//++ get from embed huh..
	return false
}

func (vfi *virtualFileInfo) Sys() interface{} {
	return nil
}
