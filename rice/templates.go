package main

import (
	"fmt"
	"os"
	"text/template"
)

var tmplEmbeddedBox *template.Template

func init() {
	var err error

	// parse embedded box template
	tmplEmbeddedBox, err = template.New("embeddedBox").Parse(`package {{.Package}}

import (
	"github.com/GeertJohan/go.rice"
	"time"
)

func init() {

	// define files
	{{range .Files}}{{.Identifier}} := &rice.EmbeddedFile{
		Filename:    ` + "`" + `{{.FileName}}` + "`" + `,
		FileModTime: time.Unix({{.ModTime}}, 0),
		Content:     string({{.Content | printf "%#v"}}), //++ TODO: optimize? (double allocation) or does compiler already optimize this?
	}
	{{end}}

	// define dirs
	{{range .Dirs}}{{.Identifier}} := &rice.EmbeddedDir{
		Filename:    ` + "`" + `{{.FileName}}` + "`" + `,
		DirModTime: time.Unix({{.ModTime}}, 0),
		ChildFiles:  []*rice.EmbeddedFile{
			{{range .ChildFiles}}{{.Identifier}}, // {{.FileName}}
			{{end}}
		},
	}
	{{end}}

	// link ChildDirs
	{{range .Dirs}}{{.Identifier}}.ChildDirs = []*rice.EmbeddedDir{
		{{range .ChildDirs}}{{.Identifier}}, // {{.FileName}}
		{{end}}
	}
	{{end}}

	// register embeddedBox
	rice.RegisterEmbeddedBox(` + "`" + `{{.BoxName}}` + "`" + `, &rice.EmbeddedBox{
		Name: ` + "`" + `{{.BoxName}}` + "`" + `,
		Time: time.Unix({{.UnixNow}}, 0),
		Dirs: map[string]*rice.EmbeddedDir{
			{{range .Dirs}}"{{.FileName}}": {{.Identifier}},
			{{end}}
		},
		Files: map[string]*rice.EmbeddedFile{
			{{range .Files}}"{{.FileName}}": {{.Identifier}},
			{{end}}
		},
	})
}`)
	if err != nil {
		fmt.Printf("error parsing embedded box template: %s\n", err)
		os.Exit(-1)
	}
}

type boxDataType struct {
	Package string
	BoxName string
	UnixNow int64
	Files   []*fileDataType
	Dirs    map[string]*dirDataType
}

type singleDataType struct {
	Package    string
	SingleName string
	UnixNow    int64
	File       *fileDataType
}

type fileDataType struct {
	Identifier string
	FileName   string
	Content    []byte
	ModTime    int64
}

type dirDataType struct {
	Identifier string
	FileName   string
	Content    []byte
	ModTime    int64
	ChildDirs  []*dirDataType
	ChildFiles []*fileDataType
}
