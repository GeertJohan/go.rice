package main

import (
	"code.google.com/p/go.tools/importer"
	"code.google.com/p/go.tools/oracle"
	"fmt"
	"go/build"
	"os"
)

func operationEmbed(path string) {
	importArgs := []string{
		path,
		"github.com/GeertJohan/go.rice/rice/oraclehook",
	}
	importerConfig := &importer.Config{
		Build: &build.Default,
	}
	oracleImporter := importer.New(importerConfig)
	or, err := oracle.New(oracleImporter, importArgs, nil, false)
	if err != nil {
		fmt.Printf("error making new oracle instance: %s\n", err)
		os.Exit(-1)
	}

	queryImporter := importer.New(importerConfig)
	_, err = queryImporter.LoadPackage("github.com/GeertJohan/go.rice/rice/oraclehook")
	if err != nil {
		fmt.Printf("loading package for queryImporter: %s\n", err)
		os.Exit(-1)
	}
	queryPos, err := oracle.ParseQueryPos(queryImporter, "github.com/GeertJohan/go.rice/rice/oraclehook/hook.go:#155", false)
	if err != nil {
		fmt.Printf("error making queryPos, oraclehook corrupt?? %s\n", err)
		os.Exit(-1)
	}
	res, err := or.Query("caller", queryPos)
	if err != nil {
		fmt.Printf("error querying caller: %s\n", err)
		os.Exit(-1)
	}
	_ = res
}
