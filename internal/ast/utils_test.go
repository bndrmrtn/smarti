package ast

import (
	"fmt"
	"testing"

	"github.com/smlgh/smarti/internal/lexer"
)

func Test_FnCall(t *testing.T) {
	name, args := getFuncCall(lexer.LexerToken{
		Type:  lexer.FuncCall,
		Value: "capitalize(\"martin\", 5, xd(\"idk\", 2), get(2))",
	})

	fmt.Println("NAME:", name)
	for _, arg := range args {
		fmt.Println("ARG:", arg)
	}
}
