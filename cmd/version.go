package cmd

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v", "ver"},
	Short:   "Displays Smarti's current version.",
	Run:     execVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func execVersion(cmd *cobra.Command, args []string) {
	cmd.Println("Smarti\nVersion 0.1.0-development")
}
