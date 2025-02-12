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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to the list.",
	Long: `Use the -m argument to add the one-line task description on the
command line, or leave it off to be prompted. Either way, you will be
asked to provide the relative priority of this task against a small
subset of the existing tasks. (The larger the task list, the smaller
the proportion of tasks you will be asked to rank.)`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := cliList()
		if err != nil {
			log.With(`err`, err).Fatal(`could not get task list`)
		}
		app := app.App{
			UI:   ui.CLI{},
			List: tasks,
		}
		task := cmd.Flags().Lookup(`task`).Value.String()
		if err := app.Add(task); err != nil {
			log.With(`err`, err).Fatal(`could not get task list`)
		}
		if err := tasks.Save(); err != nil {
			log.With(`err`, err).Fatal(`could not save task list`)
		}
	},
}

func init() {
	getRootCmd().AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP("task", "m", "", "the task to add. You'll be prompted if this is omitted.")
}
