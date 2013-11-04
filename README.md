## go.rice

go.rice is a [Go](http://golang.org) package that makes embedding files such as html,js,css and images easy.
The package wraps basic `os` pkg functionality. During development, opened files are read directly from disk.
Upon deployment it is easy to embed all dependent files

### What does it do?
go.rice makes working with resource files easy. It doesn't matter whether the resource is html,css,js,image or a template.
The first thing go.rice does is finding the correct absolute path for your files. Say you are executing a built go binary in your home directory, but your html/template files are located in `$GOPATH/src/yourCommand/templates`. go.rice will resolve any relative path given to `rice.Box(..)` relative to the directory of the go source file calling it.

Ofcourse, this only works when the actual source is available. Sometimes you wish to simply push a binary. For instance, in server deployment. This is where the `rice` tool comes in. The `rice` tool analyses source code and finds call's to `rice.Box(..)` and embeds the files in the found directories. For each box a `.go` source file is generated.

### Installation

Use `go get` for the package and `go install` for the tool.
```
go get github.com/GeertJohan/go.rice
go install github.com/GeertJohan/go.rice/rice
```

### Usage

Import the package: `import "github.com/GeertJohan/go.rice"`

**Serving HTTP from a rice.Box**
```go
http.Handle("/", http.FileServer(rice.FindBox("http-files")))
http.ListenAndServe(":8080", nil)
```

**Loading a template**
```go
// find/create a rice.Box
templateBox, err := rice.FindBox("example-templates")
if err != nil {
	log.Fatal(err)
}
// get file contents as string
templateString, err := templateBox.String("message.tmpl")
if err != nil {
	log.Fatal(err)
}
// parse and execute the template
tmplMessage, err := template.New("message").Parse(templateString)
if err != nil {
	log.Fatal(err)
}
tmplMessage.Execute(os.Stdout, map[string]string{"Message": "Hello, world!"})

```

### Licence

This project is licensed under a Simplified BSD license. Please read the [LICENSE file][license].


### Todo
 - rice.FindSingle() that loads and embeds a single file as oposed to a directory. It should have methods .String(), .Bytes() and .File()
 - think about MustString and MustBytes methods, which wrap String and Bytes, but panic on error and have single return value (string or []byte)
 - The rice tool uses a simple regexp to find calls to `rice.Box(..)`, this should use `go/ast` or maybe `go.tools/oracle`?

### Package documentation

You will find package documentation at [godoc.org/github.com/GeertJohan/go.rice][godoc].


 [license]: https://github.com/GeertJohan/go.rice/blob/master/LICENSE
 [godoc]: http://godoc.org/github.com/GeertJohan/go.rice