package botrunner

// markdown.go — convert standard Markdown to Telegram MarkdownV2.
//
// Supported conversions:
//
//	**bold**        →  *bold*
//	***bold+ital*** →  *_bold italic_*
//	*italic*        →  _italic_
//	_italic_        →  _italic_
//	~~strike~~      →  ~strike~
//	[text](url)     →  [text](url)
//	`code`          →  `code`
//	```block```     →  ```block```
//
// All other MarkdownV2 special characters are escaped with backslash.

import (
	"regexp"
	"strings"
)

// Characters that must be escaped in plain-text segments (Telegram MarkdownV2 spec).
var plainEscapeRe = regexp.MustCompile("([_*\\[\\]()~`>#+=|{}.!\\-\\\\])")

// mdPattern matches Markdown constructs in order of priority.
// Group indices (FindAllStringSubmatchIndex):
//
//	m[2..3]  = group 1: fenced code block  ```...```
//	m[4..5]  = group 2: inline code        `...`
//	m[6..7]  = group 3: bold+italic        ***...***
//	m[8..9]  = group 4: bold               **...**
//	m[10..11]= group 5: strikethrough      ~~...~~
//	m[12..13]= group 6: italic via *       *...*
//	m[14..15]= group 7: italic via _       _..._
//	m[16..17]= group 8: link text          [text]
//	m[18..19]= group 9: link url           (url)
var mdPattern = regexp.MustCompile(
	"(```[\\s\\S]*?```)" +
		"|(`[^`\\n]+`)" +
		"|\\*{3}(.+?)\\*{3}" +
		"|\\*{2}(.+?)\\*{2}" +
		"|~~(.+?)~~" +
		"|\\*([^*\\n]+?)\\*" +
		"|_([^_\\n]+?)_" +
		"|\\[([^\\]]*)\\]\\(([^)]*)\\)",
)

func escPlain(s string) string {
	return plainEscapeRe.ReplaceAllString(s, `\$1`)
}

func escLinkURL(url string) string {
	return strings.ReplaceAll(url, ")", `\)`)
}

// mdToTelegramV2 converts standard Markdown text to Telegram MarkdownV2.
// It also normalises literal "\n" escape sequences to real newlines.
func mdToTelegramV2(input string) string {
	// Treat literal \n as a real newline (common in copy-pasted messages).
	input = strings.ReplaceAll(input, `\n`, "\n")

	matches := mdPattern.FindAllStringSubmatchIndex(input, -1)

	var b strings.Builder
	last := 0

	for _, m := range matches {
		start, end := m[0], m[1]

		// Escape plain-text segment that precedes this match.
		b.WriteString(escPlain(input[last:start]))

		switch {
		case m[2] >= 0: // fenced code block — keep as-is
			b.WriteString(input[start:end])
		case m[4] >= 0: // inline code — keep as-is
			b.WriteString(input[start:end])
		case m[6] >= 0: // bold+italic ***...***
			b.WriteString("*_")
			b.WriteString(escPlain(input[m[6]:m[7]]))
			b.WriteString("_*")
		case m[8] >= 0: // bold **...**
			b.WriteByte('*')
			b.WriteString(escPlain(input[m[8]:m[9]]))
			b.WriteByte('*')
		case m[10] >= 0: // strikethrough ~~...~~
			b.WriteByte('~')
			b.WriteString(escPlain(input[m[10]:m[11]]))
			b.WriteByte('~')
		case m[12] >= 0: // italic *...*
			b.WriteByte('_')
			b.WriteString(escPlain(input[m[12]:m[13]]))
			b.WriteByte('_')
		case m[14] >= 0: // italic _..._
			b.WriteByte('_')
			b.WriteString(escPlain(input[m[14]:m[15]]))
			b.WriteByte('_')
		case m[16] >= 0: // link [text](url)
			b.WriteByte('[')
			b.WriteString(escPlain(input[m[16]:m[17]]))
			b.WriteString("](")
			b.WriteString(escLinkURL(input[m[18]:m[19]]))
			b.WriteByte(')')
		}

		last = end
	}

	// Escape the remaining plain-text tail.
	b.WriteString(escPlain(input[last:]))
	return b.String()
}
