package main

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestEmbedGo(t *testing.T) {
	sourceFiles := []sourceFile{
		{
			"boxes.go",
			[]byte(`package main

import (
	"github.com/GeertJohan/go.rice"
)

func main() {
	rice.MustFindBox("foo")
}
`),
		},
		{
			"foo/test1.txt",
			[]byte(`This is test 1`),
		},
		{
			"foo/test2.txt",
			[]byte(`This is test 2`),
		},
		{
			"foo/bar/test1.txt",
			[]byte(`This is test 1 in bar`),
		},
		{
			"foo/bar/baz/test1.txt",
			[]byte(`This is test 1 in bar/baz`),
		},
		{
			"foo/bar/baz/backtick`.txt",
			[]byte(`Backtick filename`),
		},
		{
			"foo/bar/baz/\"quote\".txt",
			[]byte(`double quoted filename`),
		},
		{
			"foo/bar/baz/'quote'.txt",
			[]byte(`single quoted filename`),
		},
		{
			"foo/`/`/`.txt",
			[]byte(`Backticks everywhere!`),
		},
		{
			"foo/new\nline",
			[]byte("File with newline in name. Yes, this is possible."),
		},
		{
			"foo/fast{%template%}",
			[]byte("Fasttemplate"),
		},
		{
			"foo/fast{%template",
			[]byte("Fasttemplate open"),
		},
		{
			"foo/fast%}template",
			[]byte("Fasttemplate close"),
		},
		{
			"foo/fast{%dir%}/test.txt",
			[]byte("Fasttemplate directory"),
		},
		{
			"foo/fast{%dir/test.txt",
			[]byte("Fasttemplate directory open"),
		},
		{
			"foo/fast%}dir/test.txt",
			[]byte("Fasttemplate directory close"),
		},
		{
			"foo/fast{$%template%$}",
			[]byte("Fasttemplate double escaping"),
		},
	}
	withIgnoredFiles := append(sourceFiles, sourceFile{"foo/rice-box.go", []byte("package main\nfunc init() {\n}")})
	pkg, cleanup, err := setUpTestPkg("foobar", withIgnoredFiles)
	defer cleanup()
	if err != nil {
		t.Error(err)
		return
	}

	var buffer bytes.Buffer

	err = writeBoxesGo(pkg, &buffer)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Generated file: \n%s", buffer.String())

	validateBoxFile(t, filepath.Join(pkg.Dir, "rice-box.go"), &buffer, sourceFiles)
}

func TestEmbedGoEmpty(t *testing.T) {
	sourceFiles := []sourceFile{
		{
			"boxes.go",
			[]byte(`package main

func main() {
}
`),
		},
	}
	pkg, cleanup, err := setUpTestPkg("foobar", sourceFiles)
	defer cleanup()
	if err != nil {
		t.Error(err)
		return
	}

	var buffer bytes.Buffer

	err = writeBoxesGo(pkg, &buffer)
	if err != errEmptyBox {
		t.Errorf("expected errEmptyBox, got %v", err)
		return
	}
}
