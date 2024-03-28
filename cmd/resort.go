/*
Copyright © 2024 Mobile Technologies Inc. <connect-support@mtigs.com>
All Rights Reserved
*/
package cmd

import (
	"log"

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
		l := 1
		if cmd.Flags().Lookup(`all`).Value.String() == `true` {
			l = tasks.Len()
		}
		for i := 0; i < l; i++ {
			task, err := tasks.Peek()
			if err != nil {
				log.Fatal(err.Error())
			}
			i, err := tasks.Pop()
			if err != nil {
				log.Fatal(err.Error())
			}
			err = ui.NewDecider(i).Run()
			if err != nil {
				log.Fatal(err.Error())
			}
			i, err = tasks.Push(task)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = ui.NewDecider(i).Run()
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		if err := tasks.Save(); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(resortCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resortCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	resortCmd.Flags().BoolP("all", "a", false, "re-sort the entire list")
}
