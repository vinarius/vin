package utils

import (
	"fmt"
	"os"

	"github.com/vinarius/vin/constants"
)

func CheckAwsProfileIsDeleteEnabled() {
	awsProfile, awsProfileIsSet := os.LookupEnv("AWS_PROFILE")

	if !awsProfileIsSet {
		fmt.Println("Aws profile is not set. Run 'vin setProfile'")
		os.Exit(0)
	}

	if awsProfile != constants.DELETE_ENABLED_AWS_PROFILE {
		fmt.Println("Don't be stupid. Set your profile correctly.")
		os.Exit(0)
	}
}
