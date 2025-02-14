package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
)

/*import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
)

type Order struct {
	OrderId string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Item    string  `json:"item"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("init aws config", slog.Any("error", err))
		os.Exit(1)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func handler(ctx context.Context, event json.RawMessage) error {
	var order Order
	if err := json.Unmarshal(event, &order); err != nil {
		slog.Error("unmarshal json", slog.Any("error", err))
		return err
	}

	bucketName := os.Getenv("RECEIPT_BUCKET")
	if bucketName == "" {
		slog.Error("RECEIPT_BUCKET env variable is not set")
		return fmt.Errorf("missing required environment variable: RECEIPT_BUCKET")
	}

	receiptContent := fmt.Sprintf("Order ID: %s\nAmount: $%.2f\nItem: %s", order.OrderId, order.Amount, order.Item)
	key := "receipts/" + order.OrderId + ".txt"

	if err := uploadReceipt(ctx, bucketName, key, receiptContent); err != nil {
		return err
	}
	slog.Info("Successfully processed order and stored receipt in bucket", slog.String("order", order.OrderId), slog.String("bucket", bucketName))
	return nil
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

*/

type Message struct {
	Msg string `json:"message"`
}

func handler(ctx context.Context) (string, error) {
	message := Message{Msg: "Hello,1211212asdasdasd world!!!"}
	raw, err := json.Marshal(message)
	return string(raw), err
}

func main() {
	lambda.Start(handler)
}
