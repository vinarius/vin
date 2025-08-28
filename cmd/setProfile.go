package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var setProfileCmd = &cobra.Command{
	Use:   "setProfile",
	Short: "A way to set your default aws profile",
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

	fmt.Println("selectedProfile:", selectedProfile)

	// TODO: setup config file and set profile in config
}
