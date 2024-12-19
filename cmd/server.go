package cmd

import (
	"fmt"
	"log"

	"github.com/bndrmrtn/smarti/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:     "server folder",
	Aliases: []string{"s", "serve", "http"},
	Short:   "Starts an HTTP server to serve the given folder.",
	Run:     execServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("listenAddr", "l", ":3000", "Address to listen on")
}

func execServer(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	srv, err := server.New(args[0])
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	listenAddr := cmd.Flag("listenAddr").Value.String()

	fmt.Printf("Server listening on %s\n", listenAddr)
	log.Fatal(srv.Start(listenAddr))
}
