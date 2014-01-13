package rice

import (
	"os"
)

// File abstracts file methods so the user doesn't see the difference between rice.virtualFile, rice.virtualDir and os.File
// This type implements the io.Reader, io.Seeker, io.Closer and http.File interfaces
type File struct {
	realF    *os.File
	virtualF *virtualFile
	virtualD *virtualDir
}

// Close is like (*os.File).Close()
// Visit http://golang.org/pkg/os/#File.Close for more information
func (f *File) Close() error {
	if f.virtualF != nil {
		return f.virtualF.close()
	}
	if f.virtualD != nil {
		return f.virtualD.close()
	}
	return f.realF.Close()
}

// Stat is like (*os.File).Stat()
// Visit http://golang.org/pkg/os/#File.Stat for more information
func (f *File) Stat() (os.FileInfo, error) {
	if f.virtualF != nil {
		return f.virtualF.stat()
	}
	if f.virtualD != nil {
		return f.virtualD.stat()
	}
	return f.realF.Stat()
}

// Readdir is like (*os.File).Readdir()
// Visit http://golang.org/pkg/os/#File.Readdir for more information
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if f.virtualF != nil {
		return f.virtualF.readdir(count)
	}
	if f.virtualD != nil {
		return f.virtualD.readdir(count)
	}
	return f.realF.Readdir(count)
}

// Read is like (*os.File).Read()
// Visit http://golang.org/pkg/os/#File.Read for more information
func (f *File) Read(bts []byte) (int, error) {
	if f.virtualF != nil {
		return f.virtualF.read(bts)
	}
	if f.virtualD != nil {
		return f.virtualD.read(bts)
	}
	return f.realF.Read(bts)
}

// Seek is like (*os.File).Seek()
// Visit http://golang.org/pkg/os/#File.Seek for more information
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.virtualF != nil {
		return f.virtualF.seek(offset, whence)
	}
	if f.virtualD != nil {
		return f.virtualD.seek(offset, whence)
	}
	return f.realF.Seek(offset, whence)
}
