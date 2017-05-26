package rice

import (
	"net/http"
	"strings"
)

// HTTPBox implements http.FileSystem which allows the use of Box with a http.FileServer.
//   e.g.: http.Handle("/", http.FileServer(rice.MustFindBox("http-files").HTTPBox()))
type HTTPBox struct {
	*Box
}

// HTTPBox creates a new HTTPBox from an existing Box
func (b *Box) HTTPBox() *HTTPBox {
	return &HTTPBox{b}
}

// Open returns a File using the http.File interface
func (hb *HTTPBox) Open(name string) (http.File, error) {
	prefix := "/" + hb.Box.name
	if strings.HasPrefix(name, prefix) {
		name = name[len(prefix):]
	}
	return hb.Box.Open(name)
}
