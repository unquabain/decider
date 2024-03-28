/*
Copyright © 2024 Mobile Technologies Inc. <connect-support@mtigs.com>
All Rights Reserved
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/Unquabain/decider/ui"
	"github.com/charmbracelet/huh"
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
		task, err := tasks.Peek()
		if err != nil {
			log.Fatal(err.Error())
		}
		confirm := cmd.Flags().Lookup(`silent`).Value.String() == `true`
		if !confirm {
			ctl := huh.NewConfirm().Title(task).Description(`Do you want to complete this task?`)
			if err := ctl.Run(); err != nil {
				log.Fatal(err.Error())
			}
			confirm = ctl.GetValue().(bool)
		}

		if !confirm {
			return
		}

		i, err := tasks.Pop()
		if err != nil {
			log.Fatal(err.Error())
		}
		d := ui.NewDecider(i)
		if err := d.Run(); err != nil {
			log.Fatal(err.Error())
		}
		if err := tasks.Save(); err != nil {
			log.Fatal(err.Error())
		}
		task, err = tasks.Peek()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(task)

	},
}

func init() {
	rootCmd.AddCommand(doneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	doneCmd.Flags().BoolP("silent", "q", false, "don't ask for confirmation")
}
