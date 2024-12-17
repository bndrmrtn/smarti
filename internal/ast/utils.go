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

func getFuncCall(s string) (string, []Node) {
	// Először a függvény nevét választjuk le a nyitó zárójel előtt
	funcName := ""
	parenStart := strings.Index(s, "(")
	if parenStart != -1 {
		funcName = s[:parenStart]
	} else {
		// Ha nincs nyitó zárójel, akkor nem függvényhívás
		return "", nil
	}

	// Argumentumok kinyerése a zárójelek között
	argsString := s[parenStart+1 : len(s)-1] // Levágjuk a zárójeleket

	// Argumentumok kinyerése és feldolgozása
	args := parseArguments(argsString)

	return funcName, args
}

// parseArguments feldolgozza a függvényhívás argumentumait
func parseArguments(argsString string) []Node {
	var args []Node
	// Először eltávolítjuk a felesleges szóközöket
	argsString = strings.TrimSpace(argsString)

	// Ha üres a string, akkor nincs argumentum
	if len(argsString) == 0 {
		return args
	}

	// Az argumentumokat felbontjuk az egyes elemekre
	argStrings := splitArguments(argsString)

	// Argumentumok Node-okká alakítása
	for _, arg := range argStrings {
		// Az argumentumot LexerToken-né alakítjuk
		t := lexer.LexerToken{Value: arg}

		// A getType segítségével meghatározzuk az értéket, típust és referencia állapotot
		value, contentType, isReference := getType(t)

		// Az argumentumból Node-ot készítünk
		args = append(args, Node{
			IsReference: isReference,
			Type:        contentType,
			Value:       value,
			Children:    nil, // Nincs gyerek elem, ha nem összetett argumentum
		})
	}

	return args
}

// splitArguments felbontja az argumentumokat
func splitArguments(argsString string) []string {
	var parts []string
	insideQuotes := false
	insideBrackets := false
	var currentArg strings.Builder

	for i := 0; i < len(argsString); i++ {
		char := argsString[i]

		// Kezeljük a stringekben lévő idézőjeleket
		if char == '"' || char == '\'' {
			insideQuotes = !insideQuotes
		}

		// Kezeljük a kocka zárójeleket
		if char == '[' {
			insideBrackets = true
		}
		if char == ']' {
			insideBrackets = false
		}

		// Ha nem vagyunk idézőjelekben és nem vagyunk kocka zárójelben
		// és elérünk egy vesszőt, akkor új argumentum kezdődik
		if char == ',' && !insideQuotes && !insideBrackets {
			parts = append(parts, currentArg.String())
			currentArg.Reset()
		} else {
			currentArg.WriteByte(char)
		}
	}

	// Ne hagyjuk el az utolsó elemet sem
	if currentArg.Len() > 0 {
		parts = append(parts, currentArg.String())
	}

	return parts
}
