package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all tasks that are still to be done",
	Long:  `Lists all tasks that are still to be done`,
	Run: func(cmd *cobra.Command, args []string) {
		results, err := tasks.List()
		if err != nil {
			fmt.Println("Something went wrong: ", err)
			return
		}
		fmt.Println("You have the following tasks: ")
		fmt.Println(results)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
