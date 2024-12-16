package lexer

func newLineLexer(char byte, i, l, p *int) {
	if char == '\n' {
		*i++
		*l++
		*p = 0
	}
}

func spaceLexer(char byte, i, p *int) {
	if char == ' ' || char == '\t' {
		*i++
		*p++
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
	newLineLexer(char, inx, line, pos)
	spaceLexer(char, inx, pos)
	singleLineCommentLexer(inx, pos, line, contentLength, content)
	multiLineCommentLexer(inx, line, pos, contentLength, content)
}

func matchToken(t Token, i int, length int, content string) bool {
	s := t.String()
	l := len(s)

	if i+l >= length {
		return false
	}

	return content[i:i+l] == t.String()
}

func letKeywordLexer(file string, inx, pos, line *int, contentLength int, content string) (*Node, error) {
	var node = new(Node)
	node.Token = Let

	assigned := false

	for *inx < contentLength {
		char := content[*inx]
		// Handle new lines, spaces and comments
		commonLexers(char, inx, line, pos, contentLength, content)

		if matchToken(Assign, *inx, contentLength, content) {
			if assigned {
				return nil, NewLexerErrorWithPosition(&ErrorPos{
					File: file,
					Line: *line,
					Pos:  *pos,
				}, "let keyword can only be assigned once")
			}

			*inx += 1
			*pos += 1
			assigned = true
		}

		if !assigned {
			*inx++
			*pos++
			continue
		}

		break
	}

	return node, nil
}
