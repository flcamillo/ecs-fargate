# cria o bucket
resource "aws_s3_bucket" "bucket_source" {
  bucket = "${data.aws_caller_identity.current.account_id}-${var.bucket_name}-source"
  tags = {
    Project = "${var.project_name}"
  }
}

# bloqueia o acesso publico ao bucket
resource "aws_s3_bucket_public_access_block" "bucket_source_block_public" {
  bucket                  = aws_s3_bucket.bucket_source.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# define a política para acesso ao bucket
data "aws_iam_policy_document" "bucket_source_policy_doc" {
  statement {
    sid    = "S3RootAccess"
    effect = "Allow"
    resources = [
      aws_s3_bucket.bucket_source.arn,
      "${aws_s3_bucket.bucket_source.arn}/*",
    ]
    actions = [
      "s3:*",
    ]
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
      ]
    }
  }
}

# define a politica de acesso ao bucket
resource "aws_s3_bucket_policy" "bucket_source_policy" {
  bucket = aws_s3_bucket.bucket_source.id
  policy = data.aws_iam_policy_document.bucket_source_policy_doc.json
}

# configura o bucket para enviar eventos de criação de objetos para o sns
resource "aws_s3_bucket_notification" "bucket_source_notification" {
  bucket = aws_s3_bucket.bucket_source.id
  topic {
    topic_arn = aws_sns_topic.topic.arn
    events    = ["s3:ObjectCreated:*"]
  }
}
