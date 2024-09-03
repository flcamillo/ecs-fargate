# cria o bucket
resource "aws_s3_bucket" "bucket_final" {
  bucket = "${data.aws_caller_identity.current.account_id}-${var.bucket_name}-final"
}

# bloqueia o acesso publico ao bucket
resource "aws_s3_bucket_public_access_block" "bucket_final_block_public" {
  bucket                  = aws_s3_bucket.bucket_final.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# define a política para acesso ao bucket
data "aws_iam_policy_document" "bucket_final_policy_doc" {
  statement {
    sid    = "S3RootAccess"
    effect = "Allow"
    resources = [
      aws_s3_bucket.bucket_final.arn,
      "${aws_s3_bucket.bucket_final.arn}/*",
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
resource "aws_s3_bucket_policy" "bucket_final_policy" {
  bucket = aws_s3_bucket.bucket_final.id
  policy = data.aws_iam_policy_document.bucket_final_policy_doc.json
}