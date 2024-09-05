# cria o topico sns
resource "aws_sns_topic" "topic" {
  name = var.sns_name
  tags = {
    Project = "${var.project_name}"
  }
}

# define a política de acesso ao topico sns
data "aws_iam_policy_document" "sns_policy_doc" {
  statement {
    sid       = "S3Publish"
    actions   = ["sns:Publish"]
    resources = [aws_sns_topic.topic.arn]
    principals {
      type        = "Service"
      identifiers = ["s3.amazonaws.com"]
    }
    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_s3_bucket.bucket_source.arn]
    }
  }
  statement {
    sid       = "OrgSubscribe"
    actions   = ["sns:Subscribe"]
    resources = [aws_sns_topic.topic.arn]
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    condition {
      test     = "StringEquals"
      variable = "aws:PrincipalOrgID"
      values   = [data.aws_organizations_organization.current.id]
    }
  }
}

# associa a política de acesso com o tópico
resource "aws_sns_topic_policy" "sns_policy" {
  arn    = aws_sns_topic.topic.arn
  policy = data.aws_iam_policy_document.sns_policy_doc.json
}
