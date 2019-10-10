package wbml

import "strings"

type Token struct {
	v []byte
	t int
}

const (
	PLAIN int = iota
	BOLD
	CODE
	ENDOFARTICLE
	ENDOFPARAGRAPH
	HEADING
	ITALIC
	LINKSTART
	LINKEND
	LINKTEXTSTART
	LINKTEXTEND
	PLAINNEXT
	QUOTE
	ILLEGAL
)

var Symbols map[byte]int

func InitParserGlobals() {
	Symbols = make(map[byte]int)
	Symbols['*'] = BOLD
	Symbols[':'] = CODE
	Symbols['!'] = HEADING
	Symbols[';'] = ENDOFARTICLE
	Symbols['\n'] = ENDOFPARAGRAPH
	Symbols['_'] = ITALIC
	Symbols['<'] = LINKSTART
	Symbols['>'] = LINKEND
	Symbols['('] = LINKTEXTSTART
	Symbols[')'] = LINKTEXTEND
	Symbols['\\'] = PLAINNEXT
	Symbols['#'] = QUOTE
}

func isSpecial(b byte) (ret bool) {
	return (b == '*' || b == '_' || b == '<' || b == '>' || b == '(' || b == ')' || b == '#' || b == ':' || b == '!' || b == ';')
}

// isNextPlain is for when you need to use a normally special character
// in a senctance.
func isNextPlain(b byte) (ret bool) {
	return b == '\\'
}

func isNewline(b byte) (ret bool) {
	return b == '\n'
}

func isPlain(b byte) (ret bool) {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b == '.') || (b == ',') || (b == '\'') || (b >= '0' && b <= '9') || (b == ' ')
}

func lex(inputStream *strings.Reader) []Token {
	b, err := inputStream.ReadByte()
	var tokStream []Token

	var treatAsPlain bool

	for {
		var bs []byte

		if isPlain(b) {
			var tok Token
			tok.t = PLAIN
			for isPlain(b) {
				bs = append(bs, b)
				b, _ = inputStream.ReadByte()
			}
			tok.v = bs
			tokStream = append(tokStream, tok)
		} else if isNewline(b) {
			var tok Token
			tok.t = Symbols[b]
			b, err = inputStream.ReadByte()
			if err != nil && err.Error() != "EOF" {
				warningf("ENDOFPARAGRAPH %s", err)
				break
			}
			tokStream = append(tokStream, tok)
		} else if isNextPlain(b) {
			var tok Token

			tok.t = PLAIN

			b, err = inputStream.ReadByte()
			if err != nil {
				warningf("NEXTPLAIN %s", err)
				break
			}
			bs = append(bs, b)
			tok.v = bs
			tokStream = append(tokStream, tok)
			b, _ = inputStream.ReadByte()
		} else if isSpecial(b) {
			var tok Token
			tok.t = Symbols[b]
			b, err = inputStream.ReadByte()
			if err != nil {
				warningf("SPECIAL %s", err)
				break
			}

			for isPlain(b) || isNextPlain(b) || treatAsPlain {
				if treatAsPlain == true {
					treatAsPlain = false
				}

				// NOTE:
				//      this is to include the case where a special symbol is the whole line
				//	    i.e. quote
				if isNewline(b) {
					break
				}

				if isNextPlain(b) {
					treatAsPlain = true
					b, err = inputStream.ReadByte()
					if err != nil {
						if err.Error() != "EOF" {
							warningf("PLAIN %s", err)
						}
					}
				}

				bs = append(bs, b)
				b, err = inputStream.ReadByte()
				if err != nil {
					if err.Error() != "EOF" {
						warningf("PLAIN %s", err)
					}
				}
			}

			tok.v = bs
			tokStream = append(tokStream, tok)
			b, err = inputStream.ReadByte()
			if err != nil && err.Error() != "EOF" {
				warningf("SPECIAL %s", err)
			}
		} else {
			if b == '\x00' {
				b, err = inputStream.ReadByte()
				if err != nil && err.Error() != "EOF" {
					warningf("UNWANTED %s", err)
				}
				break
			}
			warningf("UNWANTED %s", string(b))
			b, err = inputStream.ReadByte()
			if err != nil {
				break
			}
		}
	}

	return tokStream
}

func destorySlashR(input string) string {
	bs := []byte(input)
	var output []byte

	for _, b := range bs {
		if b == '\r' {
			continue
		}
		output = append(output, b)
	}
	return string(output)
}

// ParseToHtml takes in the string, tokenizing it, and builds
// an html string based on those tokens
func ParseToHtml(input string) string {
	input = destorySlashR(input)
	tokStream := lex(strings.NewReader(input))

	var tokStreamIndex int
	var end bool
	running := "<article class=\"blog_article\">"
	for !end {
		if tokStreamIndex >= 1 {
			running = running + "<article class=\"blog_article\">"
		}
		for !end {
			running = running + "\n<p class=\"blog_article_paragraph\">\n"
			for {
				if tokStreamIndex >= len(tokStream) {
					end = true
					break
				}
				tok := tokStream[tokStreamIndex]
				running = running
				if tok.t == ENDOFPARAGRAPH {
					break
				}
				if tok.t == ENDOFARTICLE {
					tokStreamIndex++
					goto ARTICLE_END
				}
				if tok.t == BOLD {
					running = running + "<b>" + string(tok.v) + "</b>"
					tokStreamIndex++
					continue
				}
				if tok.t == ITALIC {
					running = running + "<i>" + string(tok.v) + "</i>"
					tokStreamIndex++
					continue
				}
				if tok.t == QUOTE {
					running = running + "<p class=\"quote\">" + string(tok.v) + "</p>"
					tokStreamIndex++
					continue
				}
				if tok.t == CODE {
					running = running + "<pre>" + string(tok.v) + "</pre>"
					tokStreamIndex++
					continue
				}
				if tok.t == HEADING {
					running = running + "<h2 class=\"blog_article_heading\">" + string(tok.v) + "</h2>"
					tokStreamIndex++
					continue
				}
				running = running + string(tok.v)
				tokStreamIndex++
			}
			running = running + "\n</p>\n"
			tokStreamIndex++
			if tokStreamIndex > len(tokStream)-1 {
				end = true
				break
			}
			if tokStream[tokStreamIndex].t == ENDOFARTICLE {
				tokStreamIndex++
				break
			}
		}
	ARTICLE_END:
		running = running + "\n</article>\n"
	}
	return running
}
