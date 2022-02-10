package docs

import (
	_ "embed"
	"fmt"
)

var documentInfo = make(map[string]string)

func addCmdDocumentInfo(cmdName string, document string) {
	documentInfo[cmdName] = document
}

func init() {
	insertABFetchDocument()
}

func ShowCmdDocument(name string) {
	document := documentInfo[name]
	if len(document) == 0 {
		fmt.Printf("doesn't document for cmd:%s \n", name)
	} else {
		fmt.Println(document)
	}
}

//go:embed abfetch.md
var ABFetchDetailHelpString string
var ABFetch = "abfetch"
var insertABFetchDocument = func() {
	addCmdDocumentInfo(ABFetch, ABFetchDetailHelpString)
}
