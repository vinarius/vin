package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "vin",
	Short: "My tools",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("aws-profile", "", "Specify the AWS profile to use")
	viper.BindPFlag("aws_profile", rootCmd.PersistentFlags().Lookup("aws-profile"))
}

func initConfig() {
	homeDirectory, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configFileName := ".vin"
	viper.SetConfigName(configFileName)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(homeDirectory)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			filePath := viper.ConfigFileUsed()
			noConfigFileWasLoaded := filePath == ""

			if noConfigFileWasLoaded {
				filePath = filepath.Join(homeDirectory, configFileName+".toml")

				err = viper.WriteConfigAs(filePath)
				if err != nil {
					fmt.Println("something went wrong writing the new config file:", err)
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("something went wrong reading config file:", err)
			os.Exit(1)
		}

		return
	}

	awsProfileMaybe := viper.Get("aws_profile")

	if awsProfile, ok := awsProfileMaybe.(string); ok {
		os.Setenv("AWS_PROFILE", awsProfile)
	}
}
