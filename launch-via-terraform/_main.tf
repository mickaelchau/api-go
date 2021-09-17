resource "aws_security_group" "api_security_terra" {
  ingress {
    description = var.service_serving_inbound_rules.description
    from_port   = var.service_serving_inbound_rules.port
    to_port     = var.service_serving_inbound_rules.port
    protocol    = "tcp"
    cidr_blocks = var.service_serving_inbound_rules.authorized_addresses
  }

  ingress {
    description = var.ssh_inbound_rules.description
    from_port   = var.ssh_inbound_rules.port
    to_port     = var.ssh_inbound_rules.port
    protocol    = "tcp"
    cidr_blocks = var.ssh_inbound_rules.authorized_addresses
  }

  ingress {
    description = var.https_inbound_rules.description
    from_port   = var.https_inbound_rules.port
    to_port     = var.https_inbound_rules.port
    protocol    = "tcp"
    cidr_blocks = var.https_inbound_rules.authorized_addresses
  }

  egress {
    description = var.outbound_rules.description
    from_port   = var.outbound_rules.port
    to_port     = var.outbound_rules.port
    protocol    = "-1"
    cidr_blocks = var.outbound_rules.authorized_addresses
  }

  tags = {
    Name = "${var.service_name} security rules"
  }
}

resource "tls_private_key" "keys" {
  algorithm = var.keys_type
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

resource "aws_instance" "api_ec2" {
  ami             = var.instance_ami
  instance_type   = var.instance_type
  security_groups = [aws_security_group.api_security_terra.name]
  key_name        = aws_key_pair.ec2_keys.key_name
  tags = {
    Name = "${var.service_name} instance"
  }
}
