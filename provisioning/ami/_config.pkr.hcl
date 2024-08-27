packer {
  required_plugins {
    amazon-ami-management = {
      version = "~> 1.5.0"
      source  = "github.com/wata727/amazon-ami-management"
    }
    amazon = {
      source  = "github.com/hashicorp/amazon"
      version = "~> 1"
    }
  }
}
