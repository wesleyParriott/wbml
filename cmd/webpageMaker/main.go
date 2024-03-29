package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
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
	if WbmlFilePath == OutputFilePath {
		fmt.Println("The output file is the same as the input file!")
		fmt.Printf("The input file received: %s\nThe outputfile: %s\n", WbmlFilePath, OutputFilePath)
		os.Exit(1)
	}
	if !strings.Contains(WbmlFilePath, ".wbml") {
		fmt.Println("given a non wbml as input")
		fmt.Printf("file give: %s\n", WbmlFilePath)
		os.Exit(1)
	}

	wbmlFilePathSplit := strings.Split(WbmlFilePath, "/")
	wbmlFileName := wbmlFilePathSplit[len(wbmlFilePathSplit)-1]

	var beginningOfHTMLFile string
	if CssFilePath != "" {
		m := minify.New()
		m.AddFunc("text/css", css.Minify)

		cssFileContentsBytes, err := ioutil.ReadFile(CssFilePath)
		if err != nil {
			fmt.Printf("error when reading file: %s\n", err)
			os.Exit(1)
		}

		minifiedCssBytes, err := m.Bytes("text/css", cssFileContentsBytes)
		if err != nil {
			fmt.Printf("failed on minification: %s\n", err)
			os.Exit(1)
		}
		minifiedCssString := string(minifiedCssBytes)

		beginningOfHTMLFile = fmt.Sprintf("<html>\n<head>\n\t<meta charset=\"UTF-8\">\n\t<title>%s</title>\n\t<style>\n\t\t%s\n\t</style>\n</head>\n<body>\n", wbmlFileName, minifiedCssString)
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
		err = ioutil.WriteFile(OutputFilePath, TotalOutput, 0755)
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
