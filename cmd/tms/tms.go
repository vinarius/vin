package tms

import (
	"github.com/spf13/cobra"
	"github.com/vinarius/vin/v2/cmd"
)

var TmsCmd = &cobra.Command{
	Use:   "tms",
	Short: "Commands related to TMS (e.g. watching a camera's mqtt messages)",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(TmsCmd)
}
