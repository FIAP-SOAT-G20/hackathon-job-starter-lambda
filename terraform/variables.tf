# General Variables
variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "tc4"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
  default = "dev"
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "common_tags" {
  description = "Common tags to apply to all resources"
  type        = map(string)
  default = {
    Project   = "TC4"
    Terraform = "true"
  }
}

variable "log_level" {
  description = "Log level"
  type        = string
  default     = "INFO"
  validation {
    condition     = contains(["DEBUG", "INFO", "WARN", "ERROR"], var.log_level)
    error_message = "Log level must be one of: DEBUG, INFO, WARN, ERROR."
  }
}

variable "rds_db_username" {
  description = "Master username for the PostgreSQL RDS instances"
  type        = string
  default     = "postgres"
}

variable "rds_db_password" {
  description = "Master password for the PostgreSQL RDS instances"
  type        = string
  sensitive   = true
  default     = "postgres"
}

variable "sqs_queues" {
  description = "Map of SQS queue configurations"
  type = map(object({
    name                        = string
    delay_seconds               = number
    message_retention_seconds   = number
    visibility_timeout_seconds  = number
    fifo_queue                  = optional(bool)
    content_based_deduplication = optional(bool)
    tags                        = optional(map(string))
  }))
  default = {
    "order_queue" = {
      name                        = "order-status-updated"
      delay_seconds               = 0
      message_retention_seconds   = 1209600 # 14 days
      visibility_timeout_seconds  = 60
      fifo_queue                  = false
      content_based_deduplication = false
      tags = {
        Purpose = "Order Notifications"
      }
    }
    "kitchen_queue" = {
      name                        = "kitchen-status-updated"
      delay_seconds               = 0
      message_retention_seconds   = 1209600 # 14 days
      visibility_timeout_seconds  = 60
      fifo_queue                  = false
      content_based_deduplication = false
      tags = {
        Purpose = "Kitchen Notifications"
      }
    }
  }
}

# Customer Service Variables
variable "customer_project_name" {
  description = "Name of the customer service project"
  type        = string
  default     = "tc4-customer-service"
}

# Payment Service Variables
variable "payment_project_name" {
  description = "Name of the payment service project"
  type        = string
  default     = "tc4-payment-service"
}

# DynamoDB Variables
variable "dynamodb_billing_mode" {
  description = "DynamoDB billing mode (PROVISIONED or PAY_PER_REQUEST)"
  type        = string
  default     = "PAY_PER_REQUEST"
  validation {
    condition     = contains(["PROVISIONED", "PAY_PER_REQUEST"], var.dynamodb_billing_mode)
    error_message = "Billing mode must be either PROVISIONED or PAY_PER_REQUEST."
  }
}

variable "dynamodb_read_capacity" {
  description = "DynamoDB read capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "dynamodb_write_capacity" {
  description = "DynamoDB write capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "dynamodb_gsi_read_capacity" {
  description = "DynamoDB GSI read capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "dynamodb_gsi_write_capacity" {
  description = "DynamoDB GSI write capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "dynamodb_enable_pitr" {
  description = "Enable point-in-time recovery for DynamoDB table"
  type        = bool
  default     = true
}

variable "dynamodb_enable_encryption" {
  description = "Enable server-side encryption for DynamoDB table"
  type        = bool
  default     = true
}

# Payment Service DynamoDB Variables
variable "payment_dynamodb_billing_mode" {
  description = "Payment DynamoDB billing mode (PROVISIONED or PAY_PER_REQUEST)"
  type        = string
  default     = "PAY_PER_REQUEST"
  validation {
    condition     = contains(["PROVISIONED", "PAY_PER_REQUEST"], var.payment_dynamodb_billing_mode)
    error_message = "Billing mode must be either PROVISIONED or PAY_PER_REQUEST."
  }
}

variable "payment_dynamodb_read_capacity" {
  description = "Payment DynamoDB read capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "payment_dynamodb_write_capacity" {
  description = "Payment DynamoDB write capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "payment_dynamodb_gsi_read_capacity" {
  description = "Payment DynamoDB GSI read capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "payment_dynamodb_gsi_write_capacity" {
  description = "Payment DynamoDB GSI write capacity units (only used with PROVISIONED billing mode)"
  type        = number
  default     = 5
}

variable "payment_dynamodb_enable_pitr" {
  description = "Enable point-in-time recovery for Payment DynamoDB table"
  type        = bool
  default     = true
}

variable "payment_dynamodb_enable_encryption" {
  description = "Enable server-side encryption for Payment DynamoDB table"
  type        = bool
  default     = true
}

# Lambda Variables
variable "lambda_memory" {
  description = "Lambda memory in MB"
  type        = number
  default     = 512
}

variable "lambda_timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 30
}

# Application Variables
variable "jwt_secret" {
  description = "JWT secret key"
  type        = string
  sensitive   = true
  default     = "super-jwt-secret-key"
}
