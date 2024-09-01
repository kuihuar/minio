package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func makeSignature(t time.Time, payloadHash, region, service, secretAccessKey, algorithm string) string {
	// Create the string to sign
	stringToSign := strings.Join([]string{
		algorithm,
		t.Format("20060102T150405Z"),
		strings.Join([]string{
			t.Format("20060102"),
			region,
			service,
			"aws4_request",
		}, "/"),
		payloadHash,
	}, "\n")

	// Create the signing key
	hash := func(data string, key []byte) []byte {
		h := hmac.New(sha256.New, key)
		h.Write([]byte(data))
		return h.Sum(nil)
	}

	kDate := hash(t.Format("20060102"), []byte("AWS4"+secretAccessKey))
	kRegion := hash(region, kDate)
	kService := hash(service, kRegion)
	kSigning := hash("aws4_request", kService)

	// Sign the string
	signature := hex.EncodeToString(hash(stringToSign, kSigning))

	return signature
}
func generateCurlCommand(req *http.Request) string {
	var curlCommand string

	// 基础 curl 命令
	curlCommand += fmt.Sprintf("curl -X %s ", req.Method)

	// 添加请求头
	for name, values := range req.Header {
		for _, value := range values {
			curlCommand += fmt.Sprintf("-H \"%s: %s\" ", name, value)
		}
	}

	// 添加 URL
	curlCommand += fmt.Sprintf("\"%s\"", req.URL.String())

	// 添加请求体（如果有）
	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		body, _ := ioutil.ReadAll(req.Body)
		curlCommand += fmt.Sprintf(" -d '%s'", string(body))
	}

	return curlCommand
}

func main() {
	accessKey := "e2wTNktczAkH7QmBh258"
	secretKey := "mgs77ftOP5DG3jo8YWAYWERmgjwcns1JzVVecgn5"
	sessionToken := ""
	service := "s3"
	host := "192.168.1.137:9000"
	canonicalURI := "/dcloud/image/spark-3.1.2-bin-hadoop3.2.tgz"
	region := "bj"
	algorithm := "AWS4-HMAC-SHA256"
	apiGatewayURL := "http://" + host
	now := time.Now().UTC()
	signedHeaders := "host;x-amz-date"

	// Create the canonical request
	payloadHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // For GET requests, the payload is always an empty string
	canonicalRequest := strings.Join([]string{
		"GET",
		canonicalURI,
		"", // No query string
		strings.Join([]string{
			"host:" + host,
			"x-amz-date:" + now.Format("20060102T150405Z"),
		}, "\n"),
		"",
		signedHeaders,
		payloadHash,
	}, "\n")

	// Create the string to sign
	hashCanonicalRequest := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := makeSignature(now, hex.EncodeToString(hashCanonicalRequest[:]), region, service, secretKey, algorithm)

	// Create the authorization header
	credential := strings.Join([]string{
		accessKey,
		strings.Join([]string{
			now.Format("20060102"),
			region,
			service,
			"aws4_request",
		}, "/"),
	}, "/")
	authHeader := fmt.Sprintf("%s Credential=%s, SignedHeaders=%s, Signature=%s",
		algorithm, credential, signedHeaders, stringToSign)

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiGatewayURL+canonicalURI, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add headers
	req.Header.Add("Authorization", authHeader)
	req.Header.Add("X-Amz-Date", now.Format("20060102T150405Z"))
	req.Header.Add("X-Amz-Security-Token", sessionToken)
	req.Header.Add("Host", host)

	for k, v := range req.Header {
		fmt.Println("k:", k)
		fmt.Println("v:", v)
	}
	cmd := generateCurlCommand(req)
	fmt.Println("cmd:", cmd)
	return
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	localFilePath := "/Users/jianfenliu/Workspace/test/download/spark-3.1.2-bin-hadoop3.2.tgz2"
	outFile, err := os.Create(localFilePath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response Status:", resp.Status)
	println("File downloaded successfully:", localFilePath)

	// Read the response body
	// body, err := io.ReadAll(io.Reader(resp.Body))
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return
	// }

	// Print the response
	fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Body:", string(body))
}
