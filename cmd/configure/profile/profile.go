package profile

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vinarius/vin/v2/cmd/configure"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "A way to set the aws profile this program uses",
	Run:   profile,
}

func init() {
	configure.ConfigureCmd.AddCommand(profileCmd)
}

func profile(cmd *cobra.Command, args []string) {
	listProfilesCommand := exec.Command("aws", "configure", "list-profiles")

	listProfilesOutputRaw, err := listProfilesCommand.Output()
	if err != nil {
		fmt.Println("error running child process", err)
		os.Exit(1)
	}

	listProfilesOutputString := string(listProfilesOutputRaw)
	awsProfiles := strings.Fields(listProfilesOutputString)

	prompt := promptui.Select{
		Label: "Select Profile",
		Items: awsProfiles,
	}

	_, selectedProfile, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	viper.Set("aws_profile", selectedProfile)

	err = viper.WriteConfig()
	if err != nil {
		fmt.Println("something went wrong writing the config file", err)
	}
}
