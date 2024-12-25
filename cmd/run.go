package cmd

import (
	"os"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/lexer"
	"github.com/bndrmrtn/smarti/internal/runtime"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var runCmd = &cobra.Command{
	Use:     "run filename.smt",
	Aliases: []string{"r", "interpret", "execute"},
	Short:   "Interpret and execute .smt files",
	Run:     execRun,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("debug", "d", false, "Run the program in debug mode")
	runCmd.Flags().BoolP("color", "c", true, "Enable or disable colorized output")
}

func execRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	debug := cmd.Flag("debug").Value.String() == "true"

	// Tokenize the source code with lexer
	lx := lexer.New(args[0], args[1:]...)
	if err := lx.Parse(); err != nil {
		cmd.PrintErr(err)
		return
	}

	if debug {
		writeDebug("lexer.yaml", lx.Tokens)
	}

	// Generate abstract syntax tree from tokens
	parser := ast.NewParser(lx.Tokens)
	if err := parser.Parse(); err != nil {
		cmd.PrintErr(err)
		return
	}

	if debug {
		writeDebug("ast.yaml", parser.Nodes)
	}

	// Interpret the nodes with runtime
	runt := runtime.New()
	if err := runt.Run(args[0], parser.Nodes); err != nil {
		cmd.PrintErr(err)
		return
	}
}

func writeDebug(file string, v any) {
	f, err := os.Create(file)
	if err != nil {
		return
	}
	defer f.Close()

	_ = yaml.NewEncoder(f).Encode(v)
}
