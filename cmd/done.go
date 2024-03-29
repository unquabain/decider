//go:build !wasm
// +build !wasm

/*
Copyright © 2024 Ben C. Forsberg <benfrsbrg@gmail.com>
All Rights Reserved
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/Unquabain/decider/app"
	"github.com/Unquabain/decider/ui"
	"github.com/spf13/cobra"
)

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:     "done",
	Aliases: []string{"complete", "finish"},
	Short:   "Marks a task complete",
	Long: `Asks you to confirm the task you completed (or use the -q
option to skip this), and then asks you for the relative ranking
of some subset of the tasks in your list. The larger the task list,
the smaller the proportion of tasks you'll be asked to compare.`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := cliList()
		if err != nil {
			log.With(`err`, err).Fatal(`could not get task list`)
		}
		app := app.App{
			UI:   ui.CLI{},
			List: tasks,
		}
		silent := cmd.Flags().Lookup(`silent`).Value.String() == `true`
		if err := app.Complete(!silent); err != nil {
			log.With(`err`, err).Fatal(`could not complete task`)
		}
		if err := tasks.Save(); err != nil {
			log.With(`err`, err).Fatal(`could not save task list`)
		}
		if app.List.Len() == 0 {
			fmt.Println(`All caught up.`)
			return
		}
		task, err := app.Peek()
		if err != nil {
			log.With(`err`, err).Fatal(`could not print task`)
		}
		fmt.Println(task)

	},
}

func init() {
	getRootCmd().AddCommand(doneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	doneCmd.Flags().BoolP("silent", "q", false, "don't ask for confirmation")
}
