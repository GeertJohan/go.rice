package rice

import (
	"archive/zip"
	"bitbucket.org/kardianos/osext"
	"fmt"
	"github.com/daaku/go.zipexe"
	"strings"
)

// AppendedBox defines an appended box
type AppendedBox struct {
	Name  string               // box name
	Files map[string]*zip.File // appended files (*zip.File) by full path
}

// AppendedBoxes is a public register of appendes boxes
var AppendedBoxes = make(map[string]*AppendedBox)

func init() {
	// find if exec is appended
	thisFile, err := osext.Executable()
	if err != nil {
		return // not apended
	}
	rd, err := zipexe.Open(thisFile)
	if err != nil {
		return // not apended
	}

	for _, f := range rd.File {
		fmt.Printf("Found appended file: %s\n", f.Name)

		// get box and file name from f.Name
		fileParts := strings.SplitN(strings.TrimLeft(f.Name, "/"), "/", 2)
		boxName := fileParts[0]
		fileName := fileParts[1]

		// find box or create new one if doesn't exist
		box := AppendedBoxes[boxName]
		if box == nil {
			fmt.Printf("Creating box %s\n", boxName)
			box = &AppendedBox{
				Name:  boxName,
				Files: make(map[string]*zip.File),
			}
			AppendedBoxes[boxName] = box
		}

		// add file to box
		box.Files[fileName] = f
	}
}
