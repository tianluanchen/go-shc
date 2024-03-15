package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go-shc/shc"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of go-shc",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-shc version", shc.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
