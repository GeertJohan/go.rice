package rice

import (
	"path/filepath"
)

func (b *Box) Walk(root string, walkFn filepath.WalkFunc) error {
	if b.IsAppended() {
		//++ loop over appended data
		return nil
	}

	if b.IsEmbedded() {
		//++ loop over embedded data
		return
	}

	//++ wrap using filepath.Walk (absolutePath+root)
}
