resource "aws_sqs_queue" "queues" {
  for_each = var.sqs_queues

  name                       = each.value.name
  delay_seconds              = each.value.delay_seconds
  message_retention_seconds  = each.value.message_retention_seconds
  visibility_timeout_seconds = each.value.visibility_timeout_seconds

  fifo_queue                  = lookup(each.value, "fifo_queue", false)
  content_based_deduplication = lookup(each.value, "content_based_deduplication", false)

  tags = merge({
    Environment = var.environment
    Name        = each.value.name
  }, lookup(each.value, "tags", {}))
}

resource "aws_sqs_queue_policy" "sns_to_sqs" {
  for_each  = var.sns_topic_arn != "" ? var.sqs_queues : {}
  queue_url = aws_sqs_queue.queues[each.key].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowSNSToSendMessage"
        Effect = "Allow"
        Principal = {
          Service = "sns.amazonaws.com"
        }
        Action   = "sqs:SendMessage"
        Resource = aws_sqs_queue.queues[each.key].arn
        Condition = {
          ArnEquals = {
            "aws:SourceArn" = var.sns_topic_arn
          }
        }
      }
    ]
  })
}
