package lexer

import (
	"io"
	"os"
)

type Lexer struct {
	Tokens []LexerToken

	entryFile  string
	otherFiles []string

	pos int
}

func NewLexer(entry string, files ...string) *Lexer {
	return &Lexer{
		entryFile:  entry,
		otherFiles: files,
	}
}

func (l *Lexer) Parse() error {
	tokens, err := l.parse(l.entryFile)
	if err != nil {
		return err
	}

	l.Tokens = tokens

	return nil
}

func (l *Lexer) parse(file string) ([]LexerToken, error) {
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
		content       = string(b) + "\n" // Hozzáadunk egy \n-t a fájl végéhez, hogy biztosan befejezzük az utolsó sort
		contentLength = len(content)
		inx           = 0
		line          = 0
		pos           = 0
		tokens        []LexerToken
	)

	for inx < contentLength {
		char := content[inx]

		// Kezeljük a sortöréseket és szóközöket
		newLineLexer(char, &line, &pos)
		spaceLexer(char, &inx, &pos)

		// Kezeljük a kommenteket
		singleLineCommentLexer(&inx, &pos, &line, contentLength, content)
		multiLineCommentLexer(&inx, &line, &pos, contentLength, content)

		// Kezeljük a tokeneket
		if matchToken(Let, &inx, &line, &pos, contentLength, content) {
			tokens = append(tokens, newLexerToken(Let, nextWord(&inx, &line, &pos, contentLength, content), file, line, pos))
		}

		if matchToken(Const, &inx, &line, &pos, contentLength, content) {
			tokens = append(tokens, newLexerToken(Const, nextWord(&inx, &line, &pos, contentLength, content), file, line, pos))
		}

		if matchToken(Assign, &inx, &line, &pos, contentLength, content) {
			tokens = append(tokens, newLexerToken(Assign, nextValue(&inx, &line, &pos, contentLength, content), file, line, pos))
		}

		if isIdentifier(char) {
			funcName := nextWord(&inx, &line, &pos, contentLength, content)
			if content[inx] == '(' {
				tokens = append(tokens, functionCallLexer(funcName, &inx, &line, &pos, contentLength, content, file))
			} else if funcName != "" {
				tokens = append(tokens, newLexerToken(Identifier, funcName, file, line, pos))
			}
		}

		// Növeljük az indexet és pozíciót
		inx++
		pos++
	}

	return tokens, nil
}
