# S3 Bucket for User Avatars
resource "aws_s3_bucket" "user_avatars" {
  bucket        = "${var.project_name}-user-avatars-${var.environment}"
  force_destroy = true

  tags = {
    Name        = "${var.project_name}-user-avatars"
    Environment = var.environment
  }
}

# S3 Bucket for Blog Cover Images
resource "aws_s3_bucket" "blog_covers" {
  bucket        = "${var.project_name}-blog-covers-${var.environment}"
  force_destroy = true

  tags = {
    Name        = "${var.project_name}-blog-covers"
    Environment = var.environment
  }
}

# CORS Policy configuration (Allows frontend domain access)
resource "aws_s3_bucket_cors_configuration" "user_avatars_cors" {
  bucket = aws_s3_bucket.user_avatars.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "PUT", "POST", "DELETE"]
    allowed_origins = ["*"] # In strict production, restrict this to your exact frontend domain
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket_cors_configuration" "blog_covers_cors" {
  bucket = aws_s3_bucket.blog_covers.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "PUT", "POST", "DELETE"]
    allowed_origins = ["*"] # In strict production, restrict this to your exact frontend domain
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

# Public Access Blocks (Required to allow public uploads/views of avatars/covers)
resource "aws_s3_bucket_public_access_block" "user_avatars" {
  bucket = aws_s3_bucket.user_avatars.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_public_access_block" "blog_covers" {
  bucket = aws_s3_bucket.blog_covers.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
