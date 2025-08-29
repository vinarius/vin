package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setProfileCmd = &cobra.Command{
	Use:   "setProfile",
	Short: "A way to set the aws profile this program uses",
	Run:   setProfile,
}

func init() {
	rootCmd.AddCommand(setProfileCmd)
}

func setProfile(cmd *cobra.Command, args []string) {
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
