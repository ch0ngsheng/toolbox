package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version of toolbox",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("toolbox-1.0-beta")
	},
}

func init() {
	initFuncList = append(initFuncList, initVersion)
}

func initVersion() {
	rootCmd.AddCommand(versionCmd)
}
