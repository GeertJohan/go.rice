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

// EmbeddedFile defines an embedded file
type EmbeddedFile struct {
	Name    string // filename
	Content []byte
	ModTime time.Time
}

// RegisterEmbed registers an EmbeddedBox
func RegisterEmbed(name string, box *EmbeddedBox) {
	if _, exists := embeds[name]; exists {
		panic(fmt.Sprintf("EmbeddedBox with name `%s` exists already", name))
	}
	embeds[name] = box
}
