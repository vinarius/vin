package clean

import (
	"github.com/spf13/cobra"
	"github.com/vinarius/vin/v2/cmd"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Commands related to cleaning up old, sandbox resources (e.g. delete log groups)",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(CleanCmd)
}
