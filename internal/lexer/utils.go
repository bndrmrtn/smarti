package lexer

import (
	"fmt"
	"strings"
)

func newLineLexer(char byte, line, pos *int) {
	if char == '\n' {
		*line++
		*pos = 0
	}
}

func spaceLexer(char byte, inx, pos *int) {
	if char == ' ' || char == '\t' {
		*pos++
		*inx++
	}
}

func singleLineCommentLexer(i, p, l *int, length int, content string) {
	var (
		inx = *i
		pos = *p
	)

	if inx+1 < length && content[inx] == '/' && content[inx+1] == '/' {
		for inx < length && content[inx] != '\n' {
			inx++
			pos++
		}

		*i = inx
		*p = pos
		*l++
	}
}

func multiLineCommentLexer(i, l, p *int, length int, content string) {
	var (
		inx = *i

		line = *l
		pos  = *p
	)

	if inx+1 < length && content[inx] == '/' && content[inx+1] == '*' {
		for inx < length {
			if content[inx] == '\n' {
				line++
				pos = 0
			}

			if content[inx] == '*' && inx+1 < length && content[inx+1] == '/' {
				inx += 2
				pos += 2
				break
			}

			inx++
			pos++
		}

		*i = inx
		*l = line
		*p = pos
	}
}

func commonLexers(char byte, inx, pos, line *int, contentLength int, content string) {
	newLineLexer(char, line, pos)
	spaceLexer(char, inx, pos)
	singleLineCommentLexer(inx, pos, line, contentLength, content)
	multiLineCommentLexer(inx, line, pos, contentLength, content)
}

func isIdentifier(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_'
}

func matchToken(t Token, i, l, p *int, length int, content string) bool {
	var (
		inx  = *i
		line = *l
		pos  = *p
	)

	s := t.String()
	found := ""

	for inx+1 < length {
		if content[inx] == ' ' || content[inx] == '\t' || content[inx] == '\n' {
			if content[inx] == '\n' {
				line++
				pos = 0
			}
			inx++
			pos++
			break
		}

		found += string(content[inx])
		inx++
		pos++
	}

	if found == s {
		*i = inx
		*l = line
		*p = pos
		return true
	}

	return false
}

func nextWord(inx, line, pos *int, contentLength int, content string) string {
	var word string

	for *inx < contentLength {
		char := content[*inx]
		if char == ' ' || char == '\t' || char == '\n' {
			if char == '\n' {
				*line++
				*pos = 0
			}
			break
		}

		if !isIdentifier(char) {
			break
		}

		word += string(char)
		*inx++
		*pos++
	}

	return word
}

func nextValue(inx, line, pos *int, contentLength int, content string) string {
	var value string

	for *inx < contentLength {
		char := content[*inx]

		if string(char) == SemiColon.String() {
			return value
		}

		// Common lexers for handling basic tokenization
		commonLexers(char, inx, pos, line, contentLength, content)

		// Template processing
		if strings.HasPrefix(content[*inx:], TemplateStart.String()) {
			templateEnd := strings.Index(content[*inx:], TemplateEnd.String())
			if templateEnd == -1 {
				return value
			}

			value += content[*inx : *inx+templateEnd+len(TemplateEnd.String())]
			*pos += templateEnd + len(TemplateEnd.String())
			*inx += templateEnd + len(TemplateEnd.String())
			continue
		}

		// String literals handling
		if string(char) == SingleStringLiteral.String() || string(char) == DoubleStringLiteral.String() {
			quote := char
			value += string(char)
			*inx++
			*pos++

			for *inx < contentLength {
				char = content[*inx]

				// Handling escape sequences
				if char == '\\' && *inx+1 < contentLength {
					value += string(char) + string(content[*inx+1])
					*inx += 2
					*pos += 2
					continue
				}

				// Handling closing quote for the string
				if char == quote {
					value += string(char)
					*inx++
					*pos++
					break
				}

				// Collect characters within the string
				value += string(char)
				*inx++
				*pos++
			}

			// Skip spaces or operators after string literals
			for *inx < contentLength {
				char = content[*inx]
				if char == ' ' || char == '\n' || char == '+' {
					if char == '\n' {
						*line++
						*pos = 0
					}
					*inx++
					*pos++
					continue
				}

				return value
			}
			continue
		}

		// Default case: Collect other characters
		value += string(char)
		*inx++
		*pos++
	}

	return value
}

func nextFuncArgValue(inx, line, pos *int, contentLength int, content string) string {
	var value string

	for *inx < contentLength {
		char := content[*inx]

		// Skip spaces or newlines
		if char == ' ' || char == '\n' {
			if char == '\n' {
				*line++
				*pos = 0
			}
			*inx++
			*pos++
			continue
		}

		// If we encounter a closing parenthesis, we stop (end of arguments)
		if char == ')' {
			break
		}

		// If it's a quote (string literal), handle it like a string value
		if string(char) == SingleStringLiteral.String() || string(char) == DoubleStringLiteral.String() {
			quote := char
			value += string(char)
			*inx++
			*pos++

			for *inx < contentLength {
				char = content[*inx]

				// Handle escape sequences
				if char == '\\' && *inx+1 < contentLength {
					value += string(char) + string(content[*inx+1])
					*inx += 2
					*pos += 2
					continue
				}

				// Break on closing quote
				if char == quote {
					value += string(char)
					*inx++
					*pos++
					break
				}

				value += string(char)
				*inx++
				*pos++
			}
			continue
		}

		// Collect other characters (if it's not a special token)
		value += string(char)
		*inx++
		*pos++
	}

	return value
}

func functionCallLexer(funcName string, inx, line, pos *int, contentLength int, content, file string) LexerToken {
	var args []string

	// Move past the opening parenthesis
	*inx++
	*pos++

	// Loop to gather arguments until closing parenthesis
	for *inx < contentLength {
		char := content[*inx]
		commonLexers(char, inx, pos, line, contentLength, content)

		// Skip any spaces or newlines between arguments
		if char == ' ' || char == '\n' {
			if char == '\n' {
				*line++
				*pos = 0
			}
			*inx++
			*pos++
			continue
		}

		// Check for closing parenthesis
		if char == ')' {
			// When a closing parenthesis is encountered, break the loop
			*inx++
			*pos++
			break
		}

		// Collect argument value using the new helper function
		arg := nextFuncArgValue(inx, line, pos, contentLength, content)
		args = append(args, arg)

		// Check for comma separating arguments
		if *inx < contentLength && content[*inx] == ',' {
			*inx++
			*pos++
		}
	}

	// Skip any trailing semicolon if it exists
	if *inx < contentLength && content[*inx] == ';' {
		*inx++
		*pos++
	}

	// Return the function call token with the function name and arguments
	return newLexerToken(FuncCall, fmt.Sprintf("%s(%v)", funcName, args), file, *line, *pos)
}
