package tms

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/vinarius/vin/cmd"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch a camera's mqtt messages on tms",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			fmt.Println()
			fmt.Println("BATMAN! (╯°□°)╯︵ ┻━┻")
			fmt.Println()
		}()

		topic, _ := cmd.Flags().GetString("topic")
		generateOnly, _ := cmd.Flags().GetBool("generate-only")

		if topic == "" {
			log.Fatal("Topic cannot be empty. Use --topic flag.")
		}

		ctx := context.Background()

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("failed to load AWS config: %v", err)
		}

		iotClient := iot.NewFromConfig(cfg)

		describeEndpointOutput, err := iotClient.DescribeEndpoint(ctx, &iot.DescribeEndpointInput{
			EndpointType: aws.String("iot:Data-ATS"),
		})
		if err != nil {
			log.Fatalf("failed to describe IoT endpoint: %v", err)
		}
		iotEndpoint := *describeEndpointOutput.EndpointAddress

		creds, err := cfg.Credentials.Retrieve(ctx)
		if err != nil {
			log.Fatalf("failed to retrieve AWS credentials: %v", err)
		}

		presignedURL, err := generateSigV4URL(
			ctx,
			cfg.Region,
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
			iotEndpoint,
			time.Until(creds.Expires), // Use remaining credential lifetime for presigned URL expiry
		)
		if err != nil {
			log.Fatalf("failed to presign URL: %v", err)
		}

		if generateOnly {
			fmt.Println("Presigned URL:", presignedURL)
			fmt.Println("Generate only flag is set. Exiting without connecting.")
			return
		}

		clientID := "vin-cli-watcher-" + uuid.New().String()
		opts := mqtt.NewClientOptions()
		opts.AddBroker(presignedURL)
		opts.SetClientID(clientID)
		opts.SetOrderMatters(false) // Order doesn't matter for MQTT over WebSockets

		opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
			fmt.Printf("Connection lost: %v\n", err)
		})

		var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("TOPIC: %s\n", msg.Topic())
			fmt.Printf("MSG: %s\n", msg.Payload())
		}

		opts.SetOnConnectHandler(func(client mqtt.Client) {
			fmt.Println("Connected to MQTT broker.")
			token := client.Subscribe(topic, 0, messageHandler)
			if token.Wait() && token.Error() != nil {
				fmt.Printf("Failed to subscribe to topic '%s': %v\n", topic, token.Error())
				os.Exit(1)
			}
			fmt.Printf("Subscribed to topic: %s\n", topic)
		})

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalf("failed to connect to MQTT broker: %v", token.Error())
		}

		waitForInterrupt(client)
	},
}

func init() {
	watchCmd.Flags().String("topic", "", "The topic you want to listen to")
	watchCmd.MarkFlagRequired("topic")
	watchCmd.Flags().Bool("generate-only", false, "Generate and print the presigned URL without connecting")

	cmd.TmsCmd.AddCommand(watchCmd)
}

func waitForInterrupt(client mqtt.Client) {
	interruptChannel := make(chan os.Signal, 1)

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	<-interruptChannel

	fmt.Println("\nTermination signal received. Shutting down.")
	if client != nil && client.IsConnected() {
		client.Disconnect(250) // Disconnect with a short grace period
		fmt.Println("MQTT client disconnected.")
	}
}

func generateSigV4URL(
	_ context.Context,
	region string,
	accessKeyID string,
	secretAccessKey string,
	sessionToken string,
	endpoint string,
	expiresIn time.Duration,
) (string, error) {
	serviceName := "iotdevicegateway"
	now := time.Now().UTC()
	dateLong := now.Format("20060102T150405Z") // YYYYMMDDTHHMMSSZ
	dateShort := now.Format("20060102")        // YYYYMMDD
	host := endpoint
	canonicalHeaders := "host:" + host + "\n"
	signedHeaders := "host"
	httpMethod := "GET"
	canonicalURI := "/mqtt"

	queryParams := url.Values{}
	queryParams.Set("X-Amz-Algorithm", "AWS4-HMAC-SHA256")
	queryParams.Set("X-Amz-Credential", fmt.Sprintf("%s/%s/%s/%s/aws4_request", accessKeyID, dateShort, region, serviceName))
	queryParams.Set("X-Amz-Date", dateLong)
	queryParams.Set("X-Amz-SignedHeaders", signedHeaders)
	queryParams.Set("X-Amz-Expires", fmt.Sprintf("%d", int(expiresIn.Seconds())))

	var keys []string
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var canonicalQueryString bytes.Buffer
	for i, k := range keys {
		if i > 0 {
			canonicalQueryString.WriteString("&")
		}
		canonicalQueryString.WriteString(url.QueryEscape(k))
		canonicalQueryString.WriteString("=")
		canonicalQueryString.WriteString(url.QueryEscape(queryParams.Get(k)))
	}

	emptyStringHash := sha256Hash("")
	canonicalRequest := strings.Join([]string{
		httpMethod,
		canonicalURI,
		canonicalQueryString.String(),
		canonicalHeaders,
		signedHeaders,
		emptyStringHash,
	}, "\n")

	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		dateLong,
		fmt.Sprintf("%s/%s/%s/aws4_request", dateShort, region, serviceName),
		sha256Hash(canonicalRequest),
	}, "\n")

	kSecret := []byte("AWS4" + secretAccessKey)
	kDate := hmacSHA256(kSecret, []byte(dateShort))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(serviceName))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	signature := hex.EncodeToString(hmacSHA256(kSigning, []byte(stringToSign)))
	presignedURL := fmt.Sprintf("wss://%s%s?%s&X-Amz-Signature=%s",
		host,
		canonicalURI,
		canonicalQueryString.String(),
		signature,
	)

	if sessionToken != "" {
		presignedURL = fmt.Sprintf("%s&X-Amz-Security-Token=%s", presignedURL, url.QueryEscape(sessionToken))
	}

	return presignedURL, nil
}

func sha256Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func hmacSHA256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
