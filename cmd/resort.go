//go:build !wasm
// +build !wasm

/*
Copyright © 2024 Ben C. Forsberg <benfrsbrg@gmail.com>
All Rights Reserved
*/
package cmd

import (
	"github.com/charmbracelet/log"

	"github.com/Unquabain/decider/app"
	"github.com/Unquabain/decider/ui"
	"github.com/spf13/cobra"
)

// resortCmd represents the resort command
var resortCmd = &cobra.Command{
	Use:     "resort",
	Aliases: []string{"demote", "defer", "later"},
	Short:   "Re-sorts the list.",
	Long: `If the next task in the list is not the next task
you must do, this command asks you to re-rank it among the
other tasks in the list.

Optionally, with the -a option, it will re-rank the entire
list based on your priority responses.`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := cliList()
		if err != nil {
			log.With(`err`, err).Fatal(`could not get task list`)
		}
		app := app.App{
			UI:   ui.CLI{},
			List: tasks,
		}
		resortMeth := app.Resort
		if cmd.Flags().Lookup(`all`).Value.String() == `true` {
			resortMeth = app.ResortAll
		}
		if err := resortMeth(); err != nil {
			log.With(`err`, err).Fatal(`could not re-sort task list`)
		}
		if err := tasks.Save(); err != nil {
			log.With(`err`, err).Fatal(`could not save task list`)
		}
	},
}

func init() {
	getRootCmd().AddCommand(resortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resortCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	resortCmd.Flags().BoolP("all", "a", false, "re-sort the entire list")
}
