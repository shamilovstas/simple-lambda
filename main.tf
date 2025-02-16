terraform {
  backend "s3" {
    bucket         = "simple-lambda-state-k2i5b8al201"
    key            = "state/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "simple-lambda-locks"
    encrypt        = true
  }
}

locals {
  filename = "${path.module}/payload.zip"
}

data "archive_file" "lambda-payload" {
  type        = "zip"
  source_file = "${path.module}/bootstrap"
  output_path = "${path.module}/payload.zip"
}

resource "aws_s3_bucket" "bucket" {
  bucket        = "testingapp-report-bucket-k2i5b8al201"
  force_destroy = true
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda_iam" {
  name               = "lambda-iam"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "allow_write_s3" {
  statement {
    effect    = "Allow"
    actions   = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.bucket.arn}/receipts/*"]
  }
}

resource "aws_iam_role_policy" "allow_write_s3" {
  role   = aws_iam_role.lambda_iam.id
  policy = data.aws_iam_policy_document.allow_write_s3.json
}

resource "aws_apigatewayv2_api" "simple-lambda" {
  name          = "simple-lambda-gw"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "simple-lambda" {
  api_id      = aws_apigatewayv2_api.simple-lambda.id
  name        = "dev"
  auto_deploy = true
}

resource "aws_apigatewayv2_integration" "simple-lambda-integration" {
  api_id             = aws_apigatewayv2_api.simple-lambda.id
  integration_uri    = aws_lambda_function.simple-lambda.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  connection_type    = "INTERNET"
}

resource "aws_apigatewayv2_route" "simple-lambda" {
  api_id    = aws_apigatewayv2_api.simple-lambda.id
  route_key = "POST /hello"
  target    = "integrations/${aws_apigatewayv2_integration.simple-lambda-integration.id}"
}

resource "aws_lambda_permission" "api-gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.simple-lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.simple-lambda.execution_arn}/*/*"
}

resource "aws_lambda_function" "simple-lambda" {
  function_name = "simple-lambda-with-go"
  filename      = data.archive_file.lambda-payload.output_path
  role          = aws_iam_role.lambda_iam.arn
  handler       = "bootstrap"

  source_code_hash = data.archive_file.lambda-payload.output_base64sha256
  runtime          = "provided.al2023"

  environment {
    variables = {
      RECEIPT_BUCKET = aws_s3_bucket.bucket.id
    }
  }
}

output "base_url" {
  value = aws_apigatewayv2_stage.simple-lambda.invoke_url
}
