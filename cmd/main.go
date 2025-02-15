package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

type Message struct {
	Msg string `json:"message"`
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       json.RawMessage   `json:"body"`
}

func handler(ctx context.Context) (Response, error) {
	message := Message{Msg: "Hello,1211212asdasdasd world!!!"}
	raw, err := json.Marshal(message)

	response := Response{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       raw,
	}
	return response, err
}

func main() {
	lambda.Start(handler)
}
