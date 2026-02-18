package stacks

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfnTypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/spf13/cobra"
	"github.com/vinarius/vin/v2/cmd/clean"
	"github.com/vinarius/vin/v2/utils"
)

const (
	CHANNEL_BUFFER_SIZE = 500
	NUM_MAX_WORKERS     = 10
)

var (
	stacksCmd = &cobra.Command{
		Use:   "stacks",
		Short: "Delete your cloudformation stacks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vin clean stacks is now running...")

			utils.CheckAwsProfileIsDeleteEnabled()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			config, err := config.LoadDefaultConfig(ctx,
				config.WithRetryer(func() aws.Retryer {
					return retry.NewStandard(func(o *retry.StandardOptions) {
						o.MaxAttempts = 10
						o.MaxBackoff = 30 * time.Second
					})
				}),
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			cloudformationClient := cloudformation.NewFromConfig(config)

			listStacksPaginator := cloudformation.NewListStacksPaginator(cloudformationClient, &cloudformation.ListStacksInput{
				StackStatusFilter: []cfnTypes.StackStatus{
					cfnTypes.StackStatusCreateComplete,
					cfnTypes.StackStatusCreateFailed,
					cfnTypes.StackStatusCreateInProgress,
					cfnTypes.StackStatusDeleteFailed,
					cfnTypes.StackStatusReviewInProgress,
					cfnTypes.StackStatusRollbackInProgress,
					cfnTypes.StackStatusRollbackFailed,
					cfnTypes.StackStatusRollbackComplete,
					cfnTypes.StackStatusUpdateInProgress,
					cfnTypes.StackStatusUpdateCompleteCleanupInProgress,
					cfnTypes.StackStatusUpdateComplete,
					cfnTypes.StackStatusUpdateFailed,
					cfnTypes.StackStatusUpdateRollbackInProgress,
					cfnTypes.StackStatusUpdateRollbackFailed,
					cfnTypes.StackStatusUpdateRollbackCompleteCleanupInProgress,
					cfnTypes.StackStatusUpdateRollbackComplete,
					cfnTypes.StackStatusImportInProgress,
					cfnTypes.StackStatusImportComplete,
					cfnTypes.StackStatusImportRollbackInProgress,
					cfnTypes.StackStatusImportRollbackFailed,
					cfnTypes.StackStatusImportRollbackComplete,
				},
			})

			stackChan := make(chan *string, CHANNEL_BUFFER_SIZE)
			var workerWg sync.WaitGroup

			for i := range numWorkers {
				workerWg.Go(func() {
					deleteStackWorker(stackChan, i, ctx, cloudformationClient)
				})
			}

			for listStacksPaginator.HasMorePages() {
				listStacksOutput, err := listStacksPaginator.NextPage(ctx)
				if err != nil {
					fmt.Fprintln(os.Stderr, "error listing stacks:", err)
					os.Exit(1)
				}

				for _, stackSummary := range listStacksOutput.StackSummaries {
					stackName := stackSummary.StackName

					if *stackName != "CDKToolkit" {
						stackChan <- stackName
					}
				}
			}

			close(stackChan)
			workerWg.Wait()

			fmt.Println("Stack clean up complete")
		},
	}

	numWorkers = min(runtime.NumCPU(), NUM_MAX_WORKERS)
)

func init() {
	clean.CleanCmd.AddCommand(stacksCmd)
}

func deleteStackWorker(in <-chan *string, workerID int, ctx context.Context, cloudformationClient *cloudformation.Client) {
	for {
		select {
		case stackName, ok := <-in:
			if !ok {
				return
			}

			fmt.Println("Deleting stack:", *stackName)

			_, err := cloudformationClient.DeleteStack(ctx, &cloudformation.DeleteStackInput{
				StackName: stackName,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting stack '%v': %v\n", stackName, err)
				continue
			}

		case <-ctx.Done():
			fmt.Println("Context cancel received. Shutting down workerID:", workerID)
			return
		}
	}
}
