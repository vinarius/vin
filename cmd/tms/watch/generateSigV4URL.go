package watch

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

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
