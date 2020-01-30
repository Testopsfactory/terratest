# ---------------------------------------------------------------------------------------------------------------------
# AWS LAMBDA TERRAFORM EXAMPLE
# See test/terraform_aws_lambda_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

provider "archive" {
  version = "1.3"
}

data "archive_file" "zip" {
  type        = "zip"
  source_dir  = "${path.module}/src"
  output_path = "/tmp/.${var.function_name}.zip"
}

resource "aws_lambda_function" "lambda" {
  filename         = data.archive_file.zip.output_path
  source_code_hash = data.archive_file.zip.output_base64sha256
  function_name    = var.function_name
  role             = aws_iam_role.lambda.arn
  handler          = "lambda"
  runtime          = "go1.x"
}

resource "aws_iam_role" "lambda" {
  name = var.function_name
  assume_role_policy = jsonencode({
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
    Version : "2012-10-17"
  })
}
