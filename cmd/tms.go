package cmd

import (
	"github.com/spf13/cobra"
)

var TmsCmd = &cobra.Command{
	Use:   "tms",
	Short: "Commands related to TMS (e.g. watching a camera's mqtt messages)",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(TmsCmd)
}
