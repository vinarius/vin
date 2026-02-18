package watch

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/guregu/dynamo/v2"
	"github.com/spf13/cobra"
	"github.com/vinarius/vin/cmd/tms"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch mqtt messages. Defaults to watching all topics",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			fmt.Println()
			fmt.Println("BATMAN! (╯°□°)╯︵ ┻━┻")
			fmt.Println()
		}()

		getFlagOutput, err := getFlag(cmd)
		if err != nil {
			cmd.Help()
			fmt.Fprintf(os.Stderr, "\n%v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()

		fmt.Println("Setting up...")

		sdkCfg, err := awsCfg.LoadDefaultConfig(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load AWS config: %v", err)
			os.Exit(1)
		}

		dynamoDb := dynamo.New(sdkCfg)

		// maybe move this to a get app config util
		stage := "prod"
		if getFlagOutput != nil && getFlagOutput.stage != "" {
			stage = getFlagOutput.stage
		}

		table := dynamoDb.Table(fmt.Sprintf("iot-ue2-stateful-table-%s", stage))
		topics := []string{"#"}

		if getFlagOutput != nil {
			switch getFlagOutput.flag {
			case Topic:
				topics = getTopicsFromTopicFlag(getFlagOutput.value)
			case CameraId:
				topics = getTopicsFromCameraIdFlag(getFlagOutput.value, table, ctx, stage)
			case ProjectId:
				topics = getTopicsFromProjectIdFlag(getFlagOutput.value, table, ctx, stage)
			}
		}

		if len(topics) == 0 {
			fmt.Println("No topics to subscribe to. Exiting.")
			return
		}

		iotClient := iot.NewFromConfig(sdkCfg)

		describeEndpointOutput, err := iotClient.DescribeEndpoint(ctx, &iot.DescribeEndpointInput{
			EndpointType: aws.String("iot:Data-ATS"),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to describe IoT endpoint: %v", err)
			os.Exit(1)
		}
		iotEndpoint := *describeEndpointOutput.EndpointAddress

		creds, err := sdkCfg.Credentials.Retrieve(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to retrieve AWS credentials: %v", err)
			os.Exit(1)
		}

		presignedURL, err := generateSigV4URL(
			ctx,
			sdkCfg.Region,
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
			iotEndpoint,
			time.Until(creds.Expires),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to presign URL: %v\n", err)
			os.Exit(1)
		}

		clientID := "vin-cli-watcher-" + uuid.New().String()
		opts := mqtt.NewClientOptions()
		opts.AddBroker(presignedURL)
		opts.SetClientID(clientID)
		opts.SetOrderMatters(false)

		// cache for client-side de-duplication, for some reason subbing to
		// topicA and topicA/# and publishing to topicA/foo sends two messages
		cache := make(map[string]time.Time)
		var mu sync.Mutex
		const cacheTTL = 1 * time.Second

		opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
			fmt.Printf("Connection lost: %v\n", err)
		})

		var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			mu.Lock()
			for k, t := range cache {
				if time.Since(t) > cacheTTL {
					delete(cache, k)
				}
			}

			key := msg.Topic() + string(msg.Payload())
			if _, found := cache[key]; found {
				mu.Unlock()
				return // Is a duplicate, already processed. Ignore it.
			}
			cache[key] = time.Now()
			mu.Unlock()

			fmt.Printf("\nTime: %s\n", time.Now())
			fmt.Printf("Topic: %s\n", msg.Topic())
			fmt.Printf("Payload: %s\n", msg.Payload())
		}

		opts.SetOnConnectHandler(func(client mqtt.Client) {
			fmt.Println("Connected to MQTT broker.")

			for _, topic := range topics {
				token := client.Subscribe(topic, 1, messageHandler)
				if token.Wait() && token.Error() != nil {
					fmt.Fprintf(os.Stderr, "Failed to subscribe to topic '%s': %v\n", topic, token.Error())
					os.Exit(1)
				}

				fmt.Printf("Subscribed to topic: %s\n", topic)
			}
		})

		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to MQTT broker: %v\n", token.Error())
			os.Exit(1)
		}

		waitForInterrupt(client)
	},
}

func init() {
	watchCmd.Flags().String(Topic.String(), "", "Comma separated topic(s) you want to listen to")
	watchCmd.Flags().String(CameraId.String(), "", "Listen to a camera's topics (comma separated camera ids)")
	watchCmd.Flags().String(ProjectId.String(), "", "Listen to a project's topics (comma separated project ids)")
	watchCmd.Flags().String(Stage.String(), "", "Specify a stage to target sandbox environments (ie to target the iot-ue2-stateful-table-mark table, specify --stage \"mark\")")
	tms.TmsCmd.AddCommand(watchCmd)
}
