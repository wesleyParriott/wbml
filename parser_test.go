package wbml

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// TODO pass in testing.T for logging
func printTokens(tokStream []Token) {
	for _, tok := range tokStream {
		var typ string
		switch tok.t {
		case PLAIN:
			typ = "PLAIN"
		case BOLD:
			typ = "BOLD"
		case ITALIC:
			typ = "ITALICS"
		case LINKSTART:
			typ = "LINKSTART"
		case LINKEND:
			typ = "LINKEND"
		case LINKTEXTSTART:
			typ = "LINKTEXTSTART"
		case LINKTEXTEND:
			typ = "LINKTEXTEND"
		case QUOTE:
			typ = "QUOTE"
		case CODE:
			typ = "CODE"
		case HEADING:
			typ = "HEADING"
		case PLAINNEXT:
			typ = "PLAINNEXT"
		case ENDOFARTICLE:
			typ = "ENDOFARTICLE"
		case ENDOFPARAGRAPH:
			typ = "ENDOFPARAGRAPH"
		}
		fmt.Printf("\t[%s]%s\n", typ, string(tok.v))
	}
	fmt.Printf("\n")
}

/*
var testPattern string = `this is some plain jane text
this is some *bold* text
this is some _italicized_ text
this is a <http://localhost:8080/view/format>(link)
# this is quoted text
: this is some code\(\)
`
*/

// TODO add links (above)
var testPattern string = `! An article heading
A paragraph with plain jane text
A paragraph with *bold* text
A paragraph with _italicized_ text
# this is quoted text
: this is some code\(\)
_you're inside a thing with a slash \?_
\!\#\:\(\)?
;ENDOFARTICLE;
!New Article
`

func TestLex(t *testing.T) {
	InitParserGlobals()
	tokStream := lex(strings.NewReader(testPattern))
	printTokens(tokStream)
}

func TestParse(t *testing.T) {
	InitParserGlobals()
	ret := ParseToHtml(testPattern)
	t.Logf("%s\n", ret)
}

func TestBadFuzzer(t *testing.T) {
	// TODO
	testSymbols := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
	~!@#$%^&*()-=_+;'":[]{}|\,.<>`
	testString := ""

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 256; i++ {
		testString = testString + string(testSymbols[rand.Intn(len(testSymbols))])
	}

	InitParserGlobals()
	ret := ParseToHtml(string(testString))

	t.Logf("%s\n", ret)
}
