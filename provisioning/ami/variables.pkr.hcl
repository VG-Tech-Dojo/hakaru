variables artifacts_commit {
  type    = string
  default = "latest"
}

data "amazon-ami" "amzn2" {
  most_recent = true
  owners      = ["137112412989"]
  region      = "ap-northeast-1"
  filters     = {
    architecture     = "x86_64"
    name             = "amzn2-ami-hvm-2.0"
    root-device-type = "ebs"
    virtualization   = "hvm"
  }
}

locals {
  image_time = formatdate("YYYYMMDDhhmmss", timestamp())
}
