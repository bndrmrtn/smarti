package ast

import (
	"strconv"
	"strings"

	"github.com/smlgh/smarti/internal/lexer"
)

func getInfo(t lexer.LexerToken) NodeFileInfo {
	return NodeFileInfo{
		File: t.Info.File,
		Pos:  t.Info.Pos,
		Line: t.Info.Line,
	}
}

func getType(t lexer.LexerToken) (string, string, bool) {
	var contentType string
	var value string
	var isReference bool

	value = t.Value

	if t.Type == lexer.Assign {
		if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") {
			if value[0] == '"' {
				contentType = "string"
			} else {
				contentType = "string_single"
			}

			value = value[1 : len(value)-1]
			value = handleEscapedString(value)
			isReference = false
		} else if _, err := strconv.Atoi(value); err == nil {
			contentType = "number"
			isReference = false
		} else if isIdentifier(value) {
			contentType = "variable"
			isReference = true
		} else {
			// Ismeretlen token típus
			contentType = "unknown"
			isReference = false
		}
	} else {
		contentType = "unknown"
		isReference = false
	}

	return value, contentType, isReference
}

func handleEscapedString(s string) string {
	escapedString := strings.ReplaceAll(s, "\\\"", "\"")
	escapedString = strings.ReplaceAll(escapedString, "\\n", "\n")
	escapedString = strings.ReplaceAll(escapedString, "\\t", "\t")

	return escapedString
}

func isIdentifier(s string) bool {
	for i := 0; i < len(s); i++ {
		char := s[i]
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char == '_') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}
