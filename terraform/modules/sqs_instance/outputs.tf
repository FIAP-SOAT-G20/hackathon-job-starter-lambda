output "sqs_queue_urls" {
  description = "Map of SQS queue names to their URLs"
  value = {
    for key, queue in aws_sqs_queue.queues :
    key => queue.id
  }
}

output "sqs_queue_arns" {
  description = "Map of SQS queue names to their ARNs"
  value = {
    for key, queue in aws_sqs_queue.queues :
    key => queue.arn
  }
}