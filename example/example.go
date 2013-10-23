package main

import (
	"encoding/hex"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
)

func main() {
	box, err := rice.FindBox("example-files")
	if err != nil {
		log.Fatalf("error opening rice.Box: %s\n", err)
	}
	spew.Dump(box)
	log.Printf("box absolue path: %s", box.Dir())

	contentString, err := box.String("file.txt")
	if err != nil {
		log.Fatalf("could not read file contents as string: %s\n", err)
	}
	log.Printf("Read some file contents as string:\n%s\n", contentString)

	contentBytes, err := box.Bytes("file.txt")
	if err != nil {
		log.Fatalf("could not read file contents as byteSlice: %s\n", err)
	}
	log.Printf("Read some file contents as byteSlice:\n%s\n", hex.Dump(contentBytes))

	file, err := box.Open("file.txt")
	if err != nil {
		log.Fatalf("could not open file: %s\n", err)
	}
	spew.Dump(file)

	http.Handle("/", http.FileServer(box))
	go http.ListenAndServe(":8123", nil)
	fmt.Printf("Serving files on :8123, press ctrl-C to exit")
	select {}
}
