package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var client = lambda.New(session.New())

type lookupResult struct {
	NumberOfAddresses int      `json:"numberOfAddresses"`
	Responses         []string `json:"responses"`
}

func executeLookups(targets []string) (string, error) {
	for _, target := range targets {
		logData := map[string]string{"target": target}
		logMsg, _ := json.Marshal(logData)
		log.Println(string(logMsg))

		addresses, err := net.LookupHost(target)
		if err != nil {
			errmsg := fmt.Sprintf("job status: failed %s", target)
			return errmsg, err
		}
		res := lookupResult{Responses: addresses, NumberOfAddresses: len(addresses)}
		jsonString, _ := json.Marshal(res)
		log.Println(string(jsonString))
	}
	return "job status: success", nil
}

func handleRequest(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	// event
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("EVENT: %s", eventJSON)
	// environment variables
	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	log.Println("ALL ENV VARS:")
	for _, element := range os.Environ() {
		log.Println(element)
	}
	// request context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("REQUEST ID: %s", lc.AwsRequestID)
	// global variable
	log.Printf("FUNCTION NAME: %s", lambdacontext.FunctionName)
	// context method
	deadline, _ := ctx.Deadline()
	log.Printf("DEADLINE: %s", deadline)
	targets := []string{"rpapi.cts.imprivata.com"}
	executeLookups(targets)
	return "FunctionCount", nil
}

func main() {
	runtime.Start(handleRequest)
}
