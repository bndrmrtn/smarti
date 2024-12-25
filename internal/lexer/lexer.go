package lexer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type Lexer struct {
	Tokens []LexerToken

	entryFile  string
	otherFiles []string

	hash string

	pos int
}

func New(entry string, files ...string) *Lexer {
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

	hash := md5.New()
	hash.Write(b)
	l.hash = hex.EncodeToString(hash.Sum(nil))

	content := string(b) + "\n"
	contentLength := len(content)
	inx, line, pos := 0, 0, 0
	var tokens []LexerToken

	for inx < contentLength {
		char := content[inx]

		// Skip whitespace
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			inx++
			pos++

			if char == '\n' {
				line++
				pos = 0
			}

			continue
		}

		// Handle common lexers
		commonLexers(char, &inx, &pos, &line, contentLength, content)

		inx++
		pos++

		if char == '<' && inx < contentLength && content[inx] == '>' {
			startPos := inx - 1
			inx++ // Skip '>'
			pos++
			templateDepth := 1

			for inx < contentLength && templateDepth > 0 {
				if content[inx] == '<' && inx+1 < contentLength && content[inx+1] == '>' {
					templateDepth++
					inx += 2
					pos += 2
				} else if content[inx] == '<' && inx+1 < contentLength && content[inx+1] == '/' && inx+2 < contentLength && content[inx+2] == '>' {
					templateDepth--
					inx += 3
					pos += 3
				} else {
					inx++
					pos++
				}
			}

			if templateDepth == 0 {
				templateToken := content[startPos:inx]
				tokens = append(tokens, newLexerToken(Template, templateToken, file, line, pos))
				continue
			}

			return nil, fmt.Errorf("unbalanced template tags at line %d, pos %d", line, pos)
		}

		// Check for function call patterns: IDENTIFIER + ( ... )
		if isIdentifierChar(char) {
			startPos := inx - 1
			for inx < contentLength && isIdentifierChar(content[inx]) {
				inx++
				pos++
			}
			identifier := content[startPos:inx]

			// Check if followed by '(' to determine if it's a function call
			if inx < contentLength && content[inx] == '(' {
				funcStart := startPos
				parensCount := 1
				inx++ // Skip '('
				pos++

				for inx < contentLength && parensCount > 0 {
					if content[inx] == '(' {
						parensCount++
					} else if content[inx] == ')' {
						parensCount--
					}
					inx++
					pos++
				}

				if parensCount == 0 {
					// Function call token
					funcCall := content[funcStart:inx]
					tokens = append(tokens, newLexerToken(FuncCall, funcCall, file, line, pos))
					continue
				}

				// Syntax error: unbalanced parentheses
				return nil, fmt.Errorf("unbalanced parentheses in function call at line %d, pos %d", line, pos)
			}

			t := isToken(identifier)
			// Regular identifier token
			tokens = append(tokens, newLexerToken(t, identifier, file, line, pos))
			continue
		}

		// Handle string literals (including escape characters)
		if char == '"' {
			startPos := inx - 1
			for inx < contentLength && content[inx] != '"' {
				// Handle escape characters
				if content[inx] == '\\' {
					inx++ // skip the escape character
				}
				inx++
				pos++
			}
			inx++ // Skip closing quote
			tokens = append(tokens, newLexerToken(DoubleStringLiteral, content[startPos:inx], file, line, pos))
			continue
		}

		if char == ';' {
			tokens = append(tokens, newLexerToken(SemiColon, ";", file, line, pos))
			continue
		}

		if char == '=' {
			if inx < contentLength && content[inx] == '=' {
				tokens = append(tokens, newLexerToken(Equal, "==", file, line, pos))
				inx++
				pos++
				continue
			}
			tokens = append(tokens, newLexerToken(Assign, "=", file, line, pos))
			continue
		}

		if char == '!' && inx < contentLength && content[inx] == '=' {
			tokens = append(tokens, newLexerToken(NotEqual, "!=", file, line, pos))
			inx++
			pos++
			continue
		}

		if char == '{' {
			tokens = append(tokens, newLexerToken(CurlyBraceStart, "{", file, line, pos))
			continue
		}

		if char == '}' {
			tokens = append(tokens, newLexerToken(CurlyBraceEnd, "}", file, line, pos))
			continue
		}

		// Handle other tokens like operators or unknowns
		tokens = append(tokens, newLexerToken(Unknown, string(char), file, line, pos))
	}

	return tokens, nil
}

func (l *Lexer) Sum() string {
	return l.hash
}

func isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_') || (c >= '0' && c <= '9') || c == '.'
}

func pref(content string, length, inx int, token Token) bool {
	var l = inx - len(token.String())

	if l < 0 || inx >= length {
		return false
	}

	return content[l:inx] == token.String()
}
