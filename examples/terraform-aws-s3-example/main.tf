# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A S3 BUCKET WITH VERSIONING ENABLED INCLUDING TAGS
# See test/terraform_aws_s3_example_test.go for how to write automated tests for this code.
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_s3_bucket" "test_bucket" {
  bucket = "${local.aws_account_id}-${var.tag_bucket_name}"
  acl    = "private"

  versioning {
    enabled = true
  }

  tags {
    Name        = "${var.tag_bucket_name}"
    Environment = "${var.tag_bucket_environment}"
  }
}
