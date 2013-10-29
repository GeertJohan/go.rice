package rice

import (
	"fmt"
	"time"
)

// EmbeddedBox defines an embedded box
type EmbeddedBox struct {
	Name  string                   // box name
	Time  time.Time                // embed time
	Files map[string]*EmbeddedFile // embedded files
}

type EmbeddedSingle struct {
	Name string        // single name
	Time time.Time     // embed time
	File *EmbeddedFile // embedded file
}

// EmbeddedFile defines an embedded file
type EmbeddedFile struct {
	Name    string // filename
	Content []byte
	ModTime time.Time
}

// RegisterEmbeddedBox registers an EmbeddedBox
func RegisterEmbeddedBox(name string, box *EmbeddedBox) {
	if _, exists := embeds[name]; exists {
		panic(fmt.Sprintf("EmbeddedBox with name `%s` exists already", name))
	}
	embeds[name] = box
}
