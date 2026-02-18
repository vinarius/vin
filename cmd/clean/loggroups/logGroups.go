package loggroups

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
	"github.com/vinarius/vin/v2/cmd/clean"
	"github.com/vinarius/vin/v2/utils"
)

const (
	CHANNEL_BUFFER_SIZE                 = 500
	NUM_MAX_WORKERS                     = 10
	TICKER_INTERVAL_DURATION_IN_SECONDS = 3
)

var (
	logGroupsCmd = &cobra.Command{
		Use:   "log-groups",
		Short: "Delete your cloudwatch log groups",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("vin clean log-groups is now running...")

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

			cloudwatchLogsClient := cloudwatchlogs.NewFromConfig(config)
			describeLogGroupsPaginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(cloudwatchLogsClient, &cloudwatchlogs.DescribeLogGroupsInput{})

			logGroupChan := make(chan *string, CHANNEL_BUFFER_SIZE)
			var workerWg sync.WaitGroup
			var deleteCounter atomic.Uint64

			go func(ctx context.Context, deleteCounter *atomic.Uint64) {
				ticker := time.NewTicker(time.Second * TICKER_INTERVAL_DURATION_IN_SECONDS)

				for {
					select {
					case <-ticker.C:
						fmt.Printf("%s: Deleted %d log groups\n", time.Now(), deleteCounter.Load())
					case <-ctx.Done():
						return
					}
				}
			}(ctx, &deleteCounter)

			for i := range numWorkers {
				workerWg.Go(func() {
					deleteLogGroupWorker(logGroupChan, i, ctx, cloudwatchLogsClient, &deleteCounter)
				})
			}

			for describeLogGroupsPaginator.HasMorePages() {
				describeLogGroupsOutput, err := describeLogGroupsPaginator.NextPage(ctx)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error describing log groups:", err)
					os.Exit(1)
				}

				for _, logGroup := range describeLogGroupsOutput.LogGroups {
					logGroupChan <- logGroup.LogGroupName
				}
			}

			close(logGroupChan)
			workerWg.Wait()

			fmt.Printf("%s: Deleted %d log groups\n", time.Now(), deleteCounter.Load())
			fmt.Println("Log group clean up complete")
		},
	}

	numWorkers = min(runtime.NumCPU(), NUM_MAX_WORKERS)
)

func init() {
	clean.CleanCmd.AddCommand(logGroupsCmd)
}

func deleteLogGroupWorker(in <-chan *string, workerID int, ctx context.Context, cloudwatchLogsClient *cloudwatchlogs.Client, deleteCounter *atomic.Uint64) {
	for {
		select {
		case logGroupName, ok := <-in:
			if !ok {
				return
			}

			_, err := cloudwatchLogsClient.DeleteLogGroup(ctx, &cloudwatchlogs.DeleteLogGroupInput{
				LogGroupName: logGroupName,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Worker %d failed to delete log group '%s': %v\n", workerID, *logGroupName, err)
				continue
			}

			deleteCounter.Add(1)

		case <-ctx.Done():
			return
		}
	}
}
