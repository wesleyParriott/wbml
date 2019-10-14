package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"wbml"
)

var (
	HelpFlag bool

	CssFilePath    string
	WbmlFilePath   string
	OutputFilePath string
)

func init() {
	flag.BoolVar(&HelpFlag, "h", false, "prints this helpful message.")

	flag.StringVar(&CssFilePath, "c", "", "Optional css filename to import. If empty no css will be included.")
	flag.StringVar(&OutputFilePath, "o", "", "Optional name of the output file. Example: ./webpageMaker -w ./cool_stuff.wbml -o ./static/index.html")
	flag.StringVar(&WbmlFilePath, "w", "", "Mandatory path of the .wbml file that the html will be rendered from.")
}

func PrintUsage() {
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()
	if HelpFlag {
		PrintUsage()
	}
	if WbmlFilePath == "" {
		fmt.Println("Need a File path to the .wbml file")
		fmt.Println("Example: $> ./webpageMaker -w index.wbml")
		fmt.Println("")
		PrintUsage()
	}

	wbmlFilePathSplit := strings.Split(WbmlFilePath, "/")
	wbmlFileName := wbmlFilePathSplit[len(wbmlFilePathSplit)-1]

	var beginningOfHTMLFile string
	if CssFilePath != "" {
		cssFileContentsBytes, err := ioutil.ReadFile(CssFilePath)
		if err != nil {
			fmt.Printf("error when reading file: %s\n", err)
			os.Exit(1)
		}
		cssFileContentsString := string(cssFileContentsBytes)

		beginningOfHTMLFile = fmt.Sprintf("<html>\n<head>\n\t<meta charset=\"UTF-8\">\n\t<title>%s</title>\n\t<style>%s</style>\n</head>\n<body>\n", wbmlFileName, cssFileContentsString)
	} else {
		beginningOfHTMLFile = fmt.Sprintf("<html>\n<head>\n\t<meta charset=\"UTF-8\">\n\t<title>%s</title>\n</head>\n<body>\n", wbmlFileName)
	}

	var endOfHTMLFile = "</body>\n</html>\n"

	inputStringByte, err := ioutil.ReadFile(WbmlFilePath)
	if err != nil {
		fmt.Printf("error when reading file: %s\n", err)
		os.Exit(1)
	}
	inputString := string(inputStringByte)

	wbml.InitParserGlobals()
	htmlOutput := wbml.ParseToHtml(inputString)

	if OutputFilePath != "" {
		TotalOutput := []byte(beginningOfHTMLFile + htmlOutput + endOfHTMLFile)
		err = ioutil.WriteFile(OutputFilePath, TotalOutput, 777)
		if err != nil {
			fmt.Printf("problem when writing to output file: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf(beginningOfHTMLFile)
		fmt.Println(htmlOutput)
		fmt.Printf(endOfHTMLFile)
	}

}
