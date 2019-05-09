data "aws_iam_policy_document" "s3_bucket_policy" {
  statement {
    effect    = "Allow"
    principals {
      identifiers = ["${local.aws_account_id}"]
      type = "AWS"
    }
    actions   = ["*"]
    resources = ["${aws_s3_bucket.test_bucket.arn}/*"]
  }

  statement {
    effect    = "Deny"
    principals {
      identifiers = ["*"]
      type = "AWS"
    }
    actions   = ["*"]
    resources = ["${aws_s3_bucket.test_bucket.arn}/*"]

    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values = [
        "false",
      ]
    }
  }
}
