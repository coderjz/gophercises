package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Adds a task with name <name>",
	Long:  `Adds a task with name <name>"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No task name provided")
			return
		}
		taskName := strings.Join(args, " ")
		tasks.Add(taskName)
		fmt.Printf("Added new task %s\n", taskName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
