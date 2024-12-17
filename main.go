package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/smlgh/smarti/internal/ast"
	"github.com/smlgh/smarti/internal/lexer"
)

func main() {
	lx := lexer.NewLexer("language/app.sml")

	if err := lx.Parse(); err != nil {
		log.Fatal(err)
	}

	writeJSON("lx.json", lx.Tokens)

	p := ast.NewParser(lx.Tokens)

	if err := p.Parse(); err != nil {
		log.Fatal(err)
	}

	writeJSON("nodes.json", p.Nodes)
}

func writeJSON(fileName string, v any) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}
