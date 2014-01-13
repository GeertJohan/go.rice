## go.rice

go.rice is a [Go](http://golang.org) package that makes working with resources such as html,js,css,images and templates very easy. During development `go.rice` will load required files directly from disk. Upon deployment it is easy to embed all dependent files using the `rice` tool, without changing the source code for your package.

### What does it do?
The first thing go.rice does is finding the correct absolute path for your resource files. Say you are executing go binary in your home directory, but your html files are located in `$GOPATH/src/webApplication/html-files`. `go.rice` will lookup the aboslute path for that directory. The only thing you have to do is include the resources using `rice.FindBox("html-files")`.

This only works when the source is available to the machine executing the binary. This is always the case when the binary was installed with `go get` or `go install`. It might happen that you wish to simply provide a binary, without source. For instance in server deployment. The `rice` tool analyses source code and finds call's to `rice.FindBox(..)` and embeds the required directories an "embeddex box". For each box a `.go` source file is generated containing the resources inside the box.

### Installation

Use `go get` for the package and `go install` for the tool.
```
go get github.com/GeertJohan/go.rice
go install github.com/GeertJohan/go.rice/rice
```

### Usage & Examples

Import the package: `import "github.com/GeertJohan/go.rice"`

**Serving a static content folder over HTTP with a rice Box**
```go
http.Handle("/", http.FileServer(rice.MustFindBox("http-files").HTTPBox()))
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


### TODO & Development
This package is not completed yet. Though it already provides working embedding, some important featuers are still missing.
 - implement Readdir() correctly on virtualDir

Less important stuff:
 - rice.FindSingle(..) that loads and embeds a single file as oposed to a complete directory. It should have methods .String(), .Bytes() and .File()
 - think about MustString and MustBytes methods, which wraps String() and Bytes(), but panics on error and have single return value (string or []byte)
 - The rice tool uses a simple regexp to find calls to `rice.Box(..)`, this should be changed to `go/ast` or maybe `go.tools/oracle`?
 - idea, os/arch dependent embeds. rice checks if embedding file has _os_arch or build flags. If box is not requested by file without build-flags, then the buildflags are applied to the embed file.

### Package documentation

You will find package documentation at [godoc.org/github.com/GeertJohan/go.rice][godoc].


 [license]: https://github.com/GeertJohan/go.rice/blob/master/LICENSE
 [godoc]: http://godoc.org/github.com/GeertJohan/go.rice