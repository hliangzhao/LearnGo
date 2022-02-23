package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

// evenCmd represents the even command
var evenCmd = &cobra.Command{
	Use:   "even",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var evenSum int
		for _, value := range args {
			intValue, _ := strconv.Atoi(value)
			if intValue%2 == 0 {
				evenSum += intValue
			}
		}
		fmt.Printf("the adding result is %d\n", evenSum)
	},
}

func init() {
	// 将evenCmd设置为addCmd的子命令
	addCmd.AddCommand(evenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// evenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// evenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
