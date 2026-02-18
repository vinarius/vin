package configure

import (
	"github.com/spf13/cobra"
	"github.com/vinarius/vin/v2/cmd"
)

var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Commands related to configuring this application",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmd.RootCmd.AddCommand(ConfigureCmd)
}
