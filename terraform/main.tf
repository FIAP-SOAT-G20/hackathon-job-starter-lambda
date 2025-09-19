terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket = "hackathon-terraform-state-job-starter"
    key    = "terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = var.common_tags
  }
}


module "sns_instance" {
  source = "./modules/sns_instance"

  aws_region     = var.aws_region
  environment    = var.environment
  sns_topic_name = "payment-approved"

  depends_on = [module.eks_instance]
}

module "sqs_instance" {
  source = "./modules/sqs_instance"

  aws_region    = var.aws_region
  environment   = var.environment
  sqs_queues    = var.sqs_queues
  sns_topic_arn = module.sns_instance.sns_topic_arn

  depends_on = [module.eks_instance, module.sns_instance]
}

# SNS to SQS Subscriptions
resource "aws_sns_topic_subscription" "kitchen_queue" {
  topic_arn = module.sns_instance.sns_topic_arn
  protocol  = "sqs"
  endpoint  = module.sqs_instance.sqs_queue_arns["kitchen_queue"]

  depends_on = [module.sns_instance, module.sqs_instance]
}

module "rds_instance" {
  source = "./modules/rds_instance"

  aws_region  = var.aws_region
  environment = var.environment
  db_username = var.rds_db_username
  db_password = var.rds_db_password
  subnet_ids  = module.eks_instance.subnet_ids
  vpc_id      = module.eks_instance.vpc_id

  depends_on = [module.eks_instance]
}

# Customer Service Components (Lambda-based)
# ECR Repository
module "ecr" {
  source = "./modules/ecr"

  project_name = var.customer_project_name
  environment  = var.environment
  common_tags  = var.common_tags
}

# DynamoDB Table for Customer Service
module "dynamodb_customer" {
  source = "./modules/dynamodb_customer"

  project_name                  = var.customer_project_name
  environment                   = var.environment
  billing_mode                  = var.dynamodb_billing_mode
  read_capacity                 = var.dynamodb_read_capacity
  write_capacity                = var.dynamodb_write_capacity
  gsi_read_capacity             = var.dynamodb_gsi_read_capacity
  gsi_write_capacity            = var.dynamodb_gsi_write_capacity
  enable_point_in_time_recovery = var.dynamodb_enable_pitr
  enable_encryption             = var.dynamodb_enable_encryption
  common_tags                   = var.common_tags
}

# DynamoDB Table for Payment Service
module "dynamodb_payments" {
  source = "./modules/dynamodb_payments"

  project_name                  = var.payment_project_name
  environment                   = var.environment
  billing_mode                  = var.payment_dynamodb_billing_mode
  read_capacity                 = var.payment_dynamodb_read_capacity
  write_capacity                = var.payment_dynamodb_write_capacity
  gsi_read_capacity             = var.payment_dynamodb_gsi_read_capacity
  gsi_write_capacity            = var.payment_dynamodb_gsi_write_capacity
  enable_point_in_time_recovery = var.payment_dynamodb_enable_pitr
  enable_encryption             = var.payment_dynamodb_enable_encryption
  common_tags                   = var.common_tags
}

# Lambda Function
module "lambda" {
  source = "./modules/lambda"

  project_name     = var.project_name
  environment      = var.environment
  lambda_image_uri = var.lambda_image_uri
  lambda_memory    = var.lambda_memory
  lambda_timeout   = var.lambda_timeout
  common_tags      = var.common_tags

  environment_variables = {
    ENVIRONMENT         = var.environment
    DYNAMODB_TABLE_NAME = module.dynamodb_payments.table_name
    DYNAMODB_REGION     = var.aws_region
    JWT_SECRET          = var.jwt_secret
    LOG_LEVEL           = var.log_level
  }

  depends_on = [module.sqs_instance, module.s3_instance]
}

# API Gateway
module "api_gateway" {
  source = "./modules/api_gateway"

  project_name      = var.customer_project_name
  environment       = var.environment
  lambda_arn        = module.lambda.lambda_arn
  lambda_invoke_arn = module.lambda.lambda_invoke_arn
  common_tags       = var.common_tags

  depends_on = [module.lambda]
}

# Lambda Permission for API Gateway
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda.lambda_function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${module.api_gateway.api_gateway_execution_arn}/*/*"

  depends_on = [module.lambda, module.api_gateway]
}
