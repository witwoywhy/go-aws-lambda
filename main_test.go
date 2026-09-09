package main

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	response, err := handler(context.Background(), events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/users",
	})
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}
	if response.StatusCode != 200 || !strings.Contains(response.Body, `"name":"Alice"`) {
		t.Fatalf("unexpected response: %#v", response)
	}
}
