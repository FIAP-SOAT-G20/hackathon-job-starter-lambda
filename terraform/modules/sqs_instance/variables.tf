variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment"
  type        = string
  default     = "dev"
}

variable "sqs_queues" {
  description = "Map of SQS queue configurations"
  type = map(object({
    name                          = string
    delay_seconds                = number
    message_retention_seconds    = number
    visibility_timeout_seconds   = number
    fifo_queue                   = optional(bool)
    content_based_deduplication  = optional(bool)
    tags                         = optional(map(string))
  }))
}

variable "sns_topic_arn" {
  description = "ARN of the SNS topic that will send messages to the SQS queues"
  type        = string
  default     = ""
}
