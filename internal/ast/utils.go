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

func getType(t lexer.LexerToken) (string, NodeType, bool) {
	var contentType NodeType
	var value string
	var isReference bool

	value = t.Value

	if value == "nil" {
		contentType = VarNil
	} else if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") {
		if value[0] == '"' {
			contentType = "string"
		} else {
			contentType = "string_single"
		}

		value = value[1 : len(value)-1]
		value = handleEscapedString(value)
	} else if strings.HasPrefix(value, lexer.TemplateStart.String()) && strings.HasSuffix(value, lexer.TemplateEnd.String()) {
		startLen := len(lexer.TemplateStart.String())
		endLen := len(lexer.TemplateEnd.String())
		value = value[startLen : len(value)-endLen]
		value = strings.TrimSpace(value)
		contentType = VarTemplate
	} else if _, err := strconv.Atoi(value); err == nil {
		contentType = VarNumber
	} else if _, err := strconv.ParseFloat(value, 64); err == nil {
		contentType = VarFloat
	} else if _, err := strconv.ParseBool(value); err == nil {
		contentType = VarBool
	} else if isIdentifier(value) {
		contentType = VarVariable
		isReference = true
	} else {
		// Unknown token
		contentType = VarUnknown
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

func getFuncCall(t lexer.LexerToken) (string, []Node) {
	funcName := ""
	parenStart := strings.Index(t.Value, "(")
	if parenStart != -1 {
		funcName = t.Value[:parenStart]
	} else {
		return "", nil
	}

	argsString := t.Value[parenStart+1 : len(t.Value)-1]

	args := parseArguments(argsString, t)

	return funcName, args
}

func parseArguments(argsString string, lt lexer.LexerToken) []Node {
	var args []Node
	argsString = strings.TrimSpace(argsString)

	if len(argsString) == 0 {
		return args
	}

	argStrings := splitArguments(argsString)

	for _, arg := range argStrings {
		t := lexer.LexerToken{Value: arg, Info: lt.Info}

		value, contentType, isReference := getType(t)

		args = append(args, Node{
			IsReference: isReference,
			Type:        contentType,
			Value:       value,
			Info:        getInfo(t),
		})
	}

	return args
}

func splitArguments(argsString string) []string {
	var parts []string
	insideQuotes := false
	insideBrackets := false
	var currentArg strings.Builder

	for i := 0; i < len(argsString); i++ {
		char := argsString[i]

		if char == '"' || char == '\'' {
			insideQuotes = !insideQuotes
		}

		if char == '[' {
			insideBrackets = true
		}
		if char == ']' {
			insideBrackets = false
		}

		if char == ',' && !insideQuotes && !insideBrackets {
			parts = append(parts, currentArg.String())
			currentArg.Reset()
		} else {
			currentArg.WriteByte(char)
		}
	}

	if currentArg.Len() > 0 {
		parts = append(parts, currentArg.String())
	}

	return parts
}
