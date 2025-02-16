package main

import (
	"context"
	"encoding/json"
	"fmt"
	"lambda-test/orders"
	"lambda-test/transport"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Headers map[string]string

var (
	ApplicationJsonHeaders = Headers{"Content-Type": "application/json"}
	s3Client               *s3.Client
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("init aws config", slog.Any("error", err))
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var order orders.Order
	if err := json.Unmarshal([]byte(event.Body), &order); err != nil {
		slog.Error("unmarshal json", slog.Any("error", err))
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnprocessableEntity,
			Body:       transport.NewMessage(err.Error()).String(),
			Headers:    ApplicationJsonHeaders,
		}
		return response, nil
	}

	bucketName := os.Getenv("RECEIPT_BUCKET")
	if bucketName == "" {
		slog.Error("RECEIPT_BUCKET env variable is not set")
		message := transport.NewMessage("RECEIPT_BUCKET env variable is not set").String()
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       message,
			Headers:    ApplicationJsonHeaders,
		}, nil
	}

	receiptContent := order.Record()
	key := fmt.Sprintf("receipts/%d.txt", order.OrderId)

	if err := uploadReceipt(ctx, bucketName, key, receiptContent); err != nil {
		message := transport.NewMessage(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       message.String(),
			Headers:    ApplicationJsonHeaders,
		}, nil
	}
	slog.Info("Successfully processed orders and stored receipt in bucket", slog.Int("orders", order.OrderId), slog.String("bucket", bucketName))
	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    ApplicationJsonHeaders,
		Body:       event.Body,
	}
	return response, nil
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
