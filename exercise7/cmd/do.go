package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do <number>",
	Short: "Completes a task, it will no longer be shown",
	Long:  `Completes the task with the given number, it will no longer be shown`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			taskNum, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Pass the task number as a single argument")
				continue
			}

			err = tasks.Do(taskNum)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Did task #%d\n", taskNum)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
