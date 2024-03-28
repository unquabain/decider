/*
Copyright © 2024 Mobile Technologies Inc. <connect-support@mtigs.com>
All Rights Reserved
*/
package cmd

import (
	"github.com/Unquabain/decider/app"
	"github.com/Unquabain/decider/ui"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all tasks.",
	Long: `Prints out the list of tasks in semi-sorted order.

Task 0 is the next task to do, and the tasks descend in
urgency generally. But there is not guaranteed to be
any strict internal ranking past 0.`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := cliList()
		if err != nil {
			log.With(`err`, err).Fatal(`could not get task list`)
		}
		app := app.App{
			UI:   ui.CLI{},
			List: tasks,
		}
		if err := app.ShowList(); err != nil {
			log.With(`err`, err).Fatal(`could not show task list`)
		}
	},
}

func init() {
	getRootCmd().AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
