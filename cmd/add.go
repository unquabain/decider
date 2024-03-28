/*
Copyright © 2024 Mobile Technologies Inc. <connect-support@mtigs.com>
All Rights Reserved
*/
package cmd

import (
	"log"

	"github.com/Unquabain/decider/ui"
	"github.com/charmbracelet/huh"
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
		var task string
		task = cmd.Flags().Lookup(`task`).Value.String()
		if task == `` {
			input := huh.NewInput().Title(`What is the new task you want to perform?`).Prompt(`> `)
			if err := input.Run(); err != nil {
				log.Fatal(err.Error())
			}
			task = input.GetValue().(string)
		}
		i, err := tasks.Push(task)
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
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().StringP("task", "m", "", "the task to add. You'll be prompted if this is omitted.")
}
