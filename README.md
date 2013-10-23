## go.rice

go.rice is a [Go](http://golang.org) package that makes embedding files such as html,js,css and images easy.
The package wraps basic `os` pkg functionality. During development, opened files are read directly from disk.
Upon deployment it is easy to embed all dependent files

### Installation

Installation is simple. Use go get:
`go get github.com/GeertJohan/go.rice`

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

### Package documentation

You will find package documentation at [godoc.org/github.com/GeertJohan/go.rice][pkgdoc].


 [license]: https://github.com/GeertJohan/go.rice/blob/master/LICENSE
 [godoc]: http://godoc.org/github.com/GeertJohan/go.rice