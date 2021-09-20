resource "tls_private_key" "keys" {
  algorithm = local.keys_type
}

module "key_pair" {
  source     = "terraform-aws-modules/key-pair/aws"
  key_name   = "created_key"
  public_key = tls_private_key.keys.public_key_openssh
}

resource "aws_key_pair" "ec2_keys" {
  key_name   = "instance_ec2_keys"
  public_key = tls_private_key.keys.public_key_openssh
}

/*resource "aws_lb" "api_load_balancer" {
  name = "API load balancer"
  internal = true
  load_balancer_type = "network"

  tags = {
    Name = "${var.service_name} load balancer"
  }
}*/

resource "aws_instance" "api_ec2" {
  ami             = local.instance_ami
  instance_type   = local.instance_type
  security_groups = [aws_security_group.api_security_terra.name]
  key_name        = aws_key_pair.ec2_keys.key_name
  tags = {
    Name = "${local.application_name} instance"
  }
}
