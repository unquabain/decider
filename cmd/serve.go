//go:build !wasm
// +build !wasm

/*
Copyright © 2024 Ben C. Forsberg <benfrsbrg@gmail.com>
All Rights Reserved
*/
package cmd

import (
	"github.com/Unquabain/decider/server"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs a web server",
	Long: `Runs the task list as a web server.

By default, storage uses the browser's localStorage.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.ListenAndServe(8899); err != nil {
			log.With(`err`, err).Fatal(`Server died unexpectedly`)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
