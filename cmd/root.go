/*
Copyright © 2024 Mobile Technologies Inc. <connect-support@mtigs.com>
All Rights Reserved
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Unquabain/decider/list"
	"github.com/spf13/cobra"
)

var tasks *list.Model

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "decider",
	Short: "Manages a to-do list.",
	Long: `Uses a standard heap-sort/priority queue algorithm
to sort a list of items.

Without any options, it prints the most urgent task in your
list. Subcommands allow you to add, remove or re-sort the
list. Each time, you will be asked the relative urgency
of a pair or trio of tasks. You must pick the most urgent
of these. You should only be asked to rank log(n) tasks,
where n is the number of tasks in the list.`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := tasks.Peek()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(-1)
		}
		fmt.Println(p)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	homedir, _ := os.UserHomeDir()
	taskfile := path.Join(homedir, ".tasks")

	rootCmd.PersistentFlags().StringP("tasks", "t", taskfile, "file with the encoded task list")
	cobra.OnInitialize(initList)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

func initList() {
	f := rootCmd.PersistentFlags().Lookup("tasks")
	fname := f.Value.String()
	tasks = list.New(fname)
	if err := tasks.Open(); err != nil {
		log.Fatal(err.Error())
	}
}
