# EKS Outputs
output "cluster_name" {
  description = "Name of the EKS cluster"
  value       = module.eks_instance.cluster_name
}

output "cluster_endpoint" {
  description = "Endpoint of the EKS cluster"
  value       = module.eks_instance.cluster_endpoint
}

output "cluster_ca_certificate" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = module.eks_instance.cluster_ca_certificate
}

output "cluster_token" {
  description = "Token for EKS cluster authentication"
  value       = module.eks_instance.cluster_token
  sensitive   = true
}

# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.eks_instance.vpc_id
}

output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = module.eks_instance.vpc_cidr
}

output "subnet_ids" {
  description = "List of subnet IDs"
  value       = module.eks_instance.subnet_ids
}

# RDS Outputs
output "rds_postgres_kitchen_endpoint" {
  description = "DNS endpoint of the kitchen RDS instance"
  value       = module.rds_instance.rds_postgres_kitchen_endpoint
}

output "rds_postgres_kitchen_db_name" {
  description = "Database name of the kitchen RDS instance"
  value       = module.rds_instance.rds_postgres_kitchen_db_name
}

output "rds_postgres_order_endpoint" {
  description = "DNS endpoint of the order RDS instance"
  value       = module.rds_instance.rds_postgres_order_endpoint
}

output "rds_postgres_order_db_name" {
  description = "Database name of the order RDS instance"
  value       = module.rds_instance.rds_postgres_order_db_name
}

output "rds_username" {
  description = "Shared username for the RDS instances"
  value       = module.rds_instance.rds_username
}

# SNS Outputs
output "sns_topic_arn" {
  description = "ARN of the created SNS topic"
  value       = module.sns_instance.sns_topic_arn
}

output "sns_topic_name" {
  description = "Name of the created SNS topic"
  value       = module.sns_instance.sns_topic_name
}

# SQS Outputs
output "sqs_queue_urls" {
  description = "Map of SQS queue names to their URLs"
  value       = module.sqs_instance.sqs_queue_urls
}

output "sqs_queue_arns" {
  description = "Map of SQS queue names to their ARNs"
  value       = module.sqs_instance.sqs_queue_arns
}

# Customer Service Outputs
# API Gateway Outputs
output "customer_api_gateway_url" {
  description = "URL of the Customer Service API Gateway"
  value       = module.api_gateway.api_gateway_url
}

output "customer_api_gateway_id" {
  description = "ID of the Customer Service API Gateway"
  value       = module.api_gateway.api_gateway_id
}

# Lambda Outputs
output "customer_lambda_function_name" {
  description = "Name of the Customer Service Lambda function"
  value       = module.lambda.lambda_function_name
}

output "customer_lambda_function_arn" {
  description = "ARN of the Customer Service Lambda function"
  value       = module.lambda.lambda_arn
}

# DynamoDB Customer Service Outputs
output "customer_dynamodb_table_name" {
  description = "Name of the Customer Service DynamoDB table"
  value       = module.dynamodb_customer.table_name
}

output "customer_dynamodb_table_arn" {
  description = "ARN of the Customer Service DynamoDB table"
  value       = module.dynamodb_customer.table_arn
}

# DynamoDB Payment Service Outputs
output "payment_dynamodb_table_name" {
  description = "Name of the Payment Service DynamoDB table"
  value       = module.dynamodb_payments.table_name
}

output "payment_dynamodb_table_arn" {
  description = "ARN of the Payment Service DynamoDB table"
  value       = module.dynamodb_payments.table_arn
}

output "payment_dynamodb_order_id_gsi" {
  description = "Name of the Payment Service order_id GSI"
  value       = module.dynamodb_payments.order_id_gsi_name
}

output "payment_dynamodb_status_gsi" {
  description = "Name of the Payment Service status GSI"
  value       = module.dynamodb_payments.status_gsi_name
}


# ECR Outputs
output "customer_ecr_repository_url" {
  description = "URL of the Customer Service ECR repository"
  value       = module.ecr.repository_url
}

# Application Endpoints
output "customer_api_endpoints" {
  description = "Customer API endpoints"
  value = {
    base_url       = module.api_gateway.api_gateway_url
    auth           = "${module.api_gateway.api_gateway_url}/auth"
    customers      = "${module.api_gateway.api_gateway_url}/customers"
    customer_by_id = "${module.api_gateway.api_gateway_url}/customers/{id}"
  }
}
