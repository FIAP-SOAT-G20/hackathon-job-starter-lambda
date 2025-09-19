# Use existing LabRole (AWS Academy)
data "aws_iam_role" "lab_role" {
  name = "LabRole"
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${var.project_name}-${var.environment}-job-starter"
  retention_in_days = var.log_retention_days

  tags = var.common_tags
}

# Lambda Function
resource "aws_lambda_function" "job_starter" {
  package_type = "Image"
  image_uri    = var.lambda_image_uri
  function_name = "${var.project_name}-${var.environment}-job-starter"
  role          = data.aws_iam_role.lab_role.arn
  memory_size   = var.lambda_memory
  timeout       = var.lambda_timeout

  environment {
    variables = var.environment_variables
  }

  depends_on = [
    aws_cloudwatch_log_group.lambda_logs,
  ]


  tags = var.common_tags
}