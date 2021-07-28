# https://www.packer.io/docs/builders/amazon/ebs
source "amazon-ebs" "hakaru" {
  ami_name        = format("hakaru - %s", local.image_time)
  ami_description = "sunrise2021 hakaru server"
  region          = "ap-northeast-1"
  ena_support     = true
  sriov_support   = true

  tags            = {
    Timestamp     = timestamp()
    SourceAMI     = source.SourceAMI
    SourceAMIName = source.SourceAMIName
  }

  instance_type               = "t3.micro"
  source_ami                  = data.amazon-ami.amzn2.id
  associate_public_ip_address = true
  iam_instance_profile        = "hakaru"
  ssh_username                = "ec2-user"
  ssh_interface               = "session_manager"
  user_data_file              = "${path.cwd}/scripts/cloud.cfg.d/99_defaults.cfg"

  security_group_filter {
    filters = {
      "tag:Name" = "hakaru"
    }
  }

  subnet_filter {
    most_free = true
    random    = false
    filters   = {
      "tag:Name" = "hakaru-public"
    }
  }

  launch_block_device_mappings {
    device_name           = "/dev/xvda"
    volume_size           = 20
    volume_type           = "gp3"
    delete_on_termination = true
  }
}

build {
  name    = "hakaru"
  sources = ["source.amazon-ebs.hakaru"]

  provisioner "shell" {
    inline = ["while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for cloud-init...'; sleep 1; done"]
  }

  provisioner "file" {
    source      = "${path.cwd}/scripts.tgz"
    destination = "/var/tmp/scripts.tgz"
  }

  provisioner "shell" {
    inline = [
      "mkdir -p /var/tmp/scripts",
      "tar xvzf /var/tmp/scripts.tgz -C /var/tmp/scripts",
      "rm -rf /var/tmp/scripts.tgz",
      "sudo make -C /var/tmp/scripts/hakaru ARTIFACTS_COMMIT=${var.artifacts_commit}"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo rm -rf /var/lib/yum && sudo yum clean all",
      "sudo rm -rf /tmp/files /home/ec2-user/files",
      "sudo rm -f /home/ec2-user/etc /home/ec2-user/.ssh/authorized_keys",
      "sudo rm -f /etc/ssh/*_key /etc/ssh/*_key.pub",
      "sudo rm -f /etc/udev/rules.d/70-persistent-net.rules"
    ]
  }
}
