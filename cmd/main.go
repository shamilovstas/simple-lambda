package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Message struct {
	Msg string `json:"message"`
}

type Headers map[string]string
type Response struct {
	StatusCode int     `json:"statusCode"`
	Headers    Headers `json:"headers"`
	Body       string  `json:"body"`
}

type Order struct {
	OrderId int     `json:"order_id"`
	Amount  float64 `json:"amount"`
	Item    string  `json:"item"`
}

var (
	s3Client *s3.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("init aws config", slog.Any("error", err))
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (*Response, error) {
	var order Order
	var response Response
	response.Headers = Headers{"Content-Type": "application/json"}

	if err := json.Unmarshal([]byte(event.Body), &order); err != nil {
		slog.Error("unmarshal json", slog.Any("error", err))
		response.StatusCode = 422
		message := Message{Msg: err.Error()}
		msg, _ := json.Marshal(message)
		response.Body = string(msg)
		return &response, err
	}

	bucketName := os.Getenv("RECEIPT_BUCKET")
	if bucketName == "" {
		slog.Error("RECEIPT_BUCKET env variable is not set")

		response.StatusCode = 401
		message := Message{Msg: "RECEIPT_BUCKET env variable is not set"}
		msg, _ := json.Marshal(message)
		response.Body = string(msg)
		return &response, errors.New("RECEIPT_BUCKET env variable is not set")
	}

	receiptContent := fmt.Sprintf("Order ID: %d\nAmount: $%.2f\nItem: %s", order.OrderId, order.Amount, order.Item)
	key := fmt.Sprintf("receipts/%d.txt", order.OrderId)

	if err := uploadReceipt(ctx, bucketName, key, receiptContent); err != nil {
		response.StatusCode = 404
		message := Message{Msg: err.Error()}
		msg, _ := json.Marshal(message)
		response.Body = string(msg)
		return &response, err
	}
	slog.Info("Successfully processed order and stored receipt in bucket", slog.Int("order", order.OrderId), slog.String("bucket", bucketName))
	response.StatusCode = 201
	response.Body = string(event.Body)
	return &response, nil
}

func uploadReceipt(ctx context.Context, bucketName, key, receiptContent string) error {
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   strings.NewReader(receiptContent),
	})

	if err != nil {
		slog.Error("Failed to upload receipt", slog.Any("error", err))
		return err
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
