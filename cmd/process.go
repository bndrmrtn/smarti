package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var procCmd = &cobra.Command{
	Use:     "process code",
	Aliases: []string{"proc"},
	Short:   "Interpret and execute inline code",
	Run:     execProc,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(procCmd)
	procCmd.Flags().BoolP("debug", "d", false, "Run the program in debug mode")
	procCmd.Flags().BoolP("color", "c", true, "Enable or disable colorized output")
}

func execProc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	code := strings.Join(args, "\n")
	tmp := os.TempDir() + "/execute.tmp.smt"
	f, err := os.Create(tmp)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	_, err = f.WriteString(code)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	execRun(cmd, []string{tmp})
}
