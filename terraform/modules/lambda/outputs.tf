output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.customer_service.function_name
}

output "lambda_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.customer_service.arn
}

output "lambda_invoke_arn" {
  description = "Invoke ARN of the Lambda function"
  value       = aws_lambda_function.customer_service.invoke_arn
}

output "lambda_role_arn" {
  description = "ARN of the Lambda IAM role"
  value       = data.aws_iam_role.lab_role.arn
}

output "lambda_qualified_arn" {
  description = "Qualified ARN of the Lambda function"
  value       = aws_lambda_function.customer_service.qualified_arn
}

output "lambda_log_group_name" {
  description = "Name of the CloudWatch log group"
  value       = aws_cloudwatch_log_group.lambda_logs.name
}