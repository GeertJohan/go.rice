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
	rice.RegisterEmbeddedBox(` + "`" + `{{.BoxName}}` + "`" + `, &rice.EmbeddedBox{
		Name: ` + "`" + `{{.BoxName}}` + "`" + `,
		Time: time.Unix({{.UnixNow}}, 0),
		Files: map[string]*rice.EmbeddedFile{
			{{range .Files}}
			"{{.FileName}}": &rice.EmbeddedFile{
				Filename:    ` + "`" + `{{.FileName}}` + "`" + `,
				FileModTime: time.Unix({{.ModTime}}, 0),
				Content:     string({{.Content | printf "%#v"}}), //++ TODO: optimize? (double allocation) or does compiler already optimize this?
			},
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
}

type singleDataType struct {
	Package    string
	SingleName string
	UnixNow    int64
	File       *fileDataType
}

type fileDataType struct {
	FileName string
	Content  []byte
	ModTime  int64
}
