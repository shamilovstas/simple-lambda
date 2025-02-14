provider "aws" {
  region = "us-east-1"
}

locals {
    filename = "${path.module}/payload.zip"
}

data "archive_file" "lambda-payload" {
  type = "zip"
  source_file = "${path.module}/bootstrap"
  output_path = "${path.module}/payload.zip"
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda_iam" {
    name = "lambda-iam"
    assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_lambda_function" "simple-lambda" {
    function_name = "simple-lambda-with-go"
    filename = data.archive_file.lambda-payload.output_path
    role = aws_iam_role.lambda_iam.arn
    handler = "bootstrap"

    source_code_hash = data.archive_file.lambda-payload.output_base64sha256
    runtime = "provided.al2023"
}

resource "aws_lambda_function_url" "simple-lambda" {
    function_name = aws_lambda_function.simple-lambda.function_name
    authorization_type = "NONE"
}


output "lambda_iam" {
    description = "Lambda IAM"
    value = data.aws_iam_policy_document.assume_role.json
}

output "lambda_url" {
    value = aws_lambda_function_url.simple-lambda.function_url
}

output "payload_hash" {
    value = data.archive_file.lambda-payload.output_base64sha256
}
