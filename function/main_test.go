package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func TestMain(t *testing.T) {
	d := time.Now().Add(2500 * time.Millisecond)
	os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "cloudwatch-alert")
	ctx, _ := context.WithDeadline(context.Background(), d)
	ctx = lambdacontext.NewContext(ctx, &lambdacontext.LambdaContext{
		AwsRequestID:       "495b12a8-xmpl-4eca-8168-160484189f99",
		InvokedFunctionArn: "arn:aws:lambda:us-east-2:123456789012:function:cloudwatch-alert",
	})
	inputJSON := ReadJSONFromFile(t, "../event.json")
	var event events.CloudWatchEvent
	err := json.Unmarshal(inputJSON, &event)
	if err != nil {
		t.Errorf("could not unmarshal event. details: %v", err)
	}
	//var inputEvent CloudWatchEvent
	result, err := handleRequest(ctx, event)
	if err != nil {
		t.Log(err)
	}
	t.Log(result)
	if !strings.Contains(result, "FunctionCount") {
		t.Errorf("Output does not contain FunctionCode.")
	}
}
func ReadJSONFromFile(t *testing.T, inputFile string) []byte {
	inputJSON, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJSON
}
