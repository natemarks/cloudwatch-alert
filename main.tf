provider "aws" {
  region = "us-east-1"
}

resource "aws_iam_role" "this" {
  name = "cloudwatch-event"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}


# -----------------------------------------------------------------------------
# AWS managed policies to lambda
# -----------------------------------------------------------------------------

resource "aws_iam_role_policy_attachment" "AWSLambdaBasicExecutionRole" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "AWSLambdaReadOnlyAccess" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AWSLambdaReadOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "AWSXrayWriteOnlyAccess" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "bin/cloudwatch-event"
  output_path = "cloudwatch-event.zip"
}

resource "aws_lambda_function" "this" {
  filename         = "cloudwatch-event.zip"
  function_name    = "cloudwatch-event"
  role             = aws_iam_role.this.arn
  handler          = "cloudwatch-event"
  memory_size      = 128
  source_code_hash = filebase64sha256("cloudwatch-event.zip")

  runtime = "go1.x"
  tracing_config {
    mode = "Active"
  }
  environment {
    variables = {
      foo = "bar"
    }
  }
}


resource "aws_cloudwatch_event_rule" "every_one_minute" {
  name                = "every-one-minute"
  description         = "Fires every one minutes"
  schedule_expression = "rate(1 minute)"
}

resource "aws_cloudwatch_event_target" "this" {
  rule      = aws_cloudwatch_event_rule.every_one_minute.name
  target_id = "lambda"
  arn       = aws_lambda_function.this.arn
}

resource "aws_lambda_permission" "this" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_one_minute.arn
}
