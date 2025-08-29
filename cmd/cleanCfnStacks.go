package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/spf13/cobra"
)

var cleanCfnStacksCmd = &cobra.Command{
	Use:   "cleanCfnStacks",
	Short: "A way to quickly delete cloudformation stacks",
	Run:   cleanCfnStacks,
}

func init() {
	rootCmd.AddCommand(cleanCfnStacksCmd)

	// TODO: add support for stack name filtering
	// TODO: add support for dry run
}

func cleanCfnStacks(cmd *cobra.Command, args []string) {
	awsProfile, awsProfileIsSet := os.LookupEnv("AWS_PROFILE")

	if !awsProfileIsSet {
		fmt.Println("Aws profile is not set. Run 'vin setProfile'")
		os.Exit(0)
	}

	if awsProfile != "t" {
		fmt.Println("Don't be stupid. Set your profile correctly.")
		os.Exit(0)
	}

	ctx := context.TODO()

	config, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal(err)
	}

	cloudformationClient := cloudformation.NewFromConfig(config)

	listStacksOutput, err := cloudformationClient.ListStacks(ctx, &cloudformation.ListStacksInput{
		StackStatusFilter: []cfnTypes.StackStatus{
			cfnTypes.StackStatusCreateComplete,
			cfnTypes.StackStatusUpdateComplete,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("listStacksOutput:")

	for i, stackSummary := range listStacksOutput.StackSummaries {
		stackName := *stackSummary.StackName

		if !strings.HasPrefix(stackName, "rec") && stackName != "CDKToolkit" {
			fmt.Printf("i: %v | stackName: %v\n", i, stackName)

			_, err := cloudformationClient.DeleteStack(ctx, &cloudformation.DeleteStackInput{
				StackName: aws.String(stackName),
			})

			if err != nil {
				fmt.Printf("Error deleting stack '%v': %v\n", stackName, err)
				continue
			}
		}
	}

	fmt.Println("cleanUpCfnStacks complete")
}
