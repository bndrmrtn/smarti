package lexer

import (
	"io"
	"os"
)

type Lexer struct {
	Nodes []Node

	entryFile  string
	otherFiles []string

	pos int
}

func NewInterpreter(entry string, files ...string) (*Lexer, error) {
	return &Lexer{
		entryFile:  entry,
		otherFiles: files,
	}, nil
}

func (l *Lexer) Parse() error {
	nodes, err := l.parse(l.entryFile)
	if err != nil {
		return err
	}

	l.Nodes = nodes

	return nil
}

func (l *Lexer) parse(file string) ([]Node, error) {
	osFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer osFile.Close()

	b, err := io.ReadAll(osFile)
	if err != nil {
		return nil, err
	}

	var (
		content       = string(b) + "\n"
		contentLength = len(content)
		inx           = 0

		line = 0
		pos  = 0
	)

	for inx < contentLength {
		char := content[inx]

		// Handle new lines and spaces
		newLineLexer(char, &inx, &line, &pos)
		spaceLexer(char, &inx, &pos)

		// Handle comments
		singleLineCommentLexer(&inx, &pos, &line, contentLength, content)
		multiLineCommentLexer(&inx, &line, &pos, contentLength, content)

		// Handle tokens
		if matchToken(Let, inx, contentLength, content) {
			letKeywordLexer(file, &inx, &pos, &line, contentLength, content)
		}

		// Increment index
		inx++
		pos++
	}

	return nil, nil
}
