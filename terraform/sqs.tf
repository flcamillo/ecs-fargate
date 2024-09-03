# cria o topico sqs
resource "aws_sqs_queue" "sqs" {
  name = var.sns_name
}

# define a política de acesso a fila sqs
data "aws_iam_policy_document" "sqs_policy_doc" {
  statement {
    sid       = "SNSSendMessage"
    effect    = "Allow"
    actions   = ["sqs:SendMessage"]
    resources = [aws_sqs_queue.sqs.arn]
    principals {
      type        = "Service"
      identifiers = ["sns.amazonaws.com"]
    }
    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_sns_topic.topic.arn]
    }
  }
}

# associa a política de acesso a fila sqs
resource "aws_sqs_queue_policy" "sqs_policy" {
  queue_url = aws_sqs_queue.sqs.id
  policy    = data.aws_iam_policy_document.sqs_policy_doc.json
}

# cria a fila sqs para mensagens mortas
resource "aws_sqs_queue" "sqs_dlq" {
  name = "${var.sns_name}_dlp"
  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = [aws_sqs_queue.sqs.arn]
  })
}

# associa a fila sqs a sua fila de mensagens mortas
resource "aws_sqs_queue_redrive_policy" "q" {
  queue_url = aws_sqs_queue.sqs.id
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.sqs_dlq.arn
    maxReceiveCount     = 4
  })
}

# cadastra a fila sqs no topico sns para receber as notificações
resource "aws_sns_topic_subscription" "subscribe_sqs_in_sns" {
  topic_arn           = aws_sns_topic.topic.arn
  protocol            = "sqs"
  endpoint            = aws_sqs_queue.sqs.arn
}
