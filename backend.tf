resource "aws_s3_bucket" "state-backend" {
  bucket = "simple-lambda-state-k2i5b8al201"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "state-backend" {
  bucket = aws_s3_bucket.state-backend.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "state-backend" {
  bucket                  = aws_s3_bucket.state-backend.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_dynamodb_table" "state-backend" {
  name         = "simple-lambda-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}