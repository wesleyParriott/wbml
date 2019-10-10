package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"wbml"
)

func main() {
	inputStringByte, err := ioutil.ReadFile("./hello.wbml")
	if err != nil {
		fmt.Printf("error when reading file: %s\n", err)
		os.Exit(1)
	}
	inputString := string(inputStringByte)

	wbml.InitParserGlobals()
	htmlOutput := wbml.ParseToHtml(inputString)
	fmt.Println(htmlOutput)
}
