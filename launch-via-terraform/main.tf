provider "aws" {
    region = "eu-west-3"
    secret_key = "X"
    access_key = "X"
}

resource "aws_security_group" "api_security_terra" {
    ingress {
        description      = "open port for the app"
        from_port = 8000
        to_port = 8000
        protocol = "tcp"
        cidr_blocks      = ["0.0.0.0/0"]
    }
    tags = {
        Name = "Terraform"
    }
    ingress {
        description      = "open ssh"
        from_port = 22
        to_port = 22
        protocol = "tcp" //has to change
        cidr_blocks      = ["0.0.0.0/0"]
    }
    ingress {
        description      = "open https"
        from_port = 443
        to_port = 443
        protocol = "tcp"
        cidr_blocks      = ["0.0.0.0/0"]
    }
    egress {
        description      = "Allow all traffic to go outside"
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks      = ["0.0.0.0/0"]
    }
}

resource "tls_private_key" "this" {
  algorithm = "RSA"
}

module "key_pair" {
  source = "terraform-aws-modules/key-pair/aws"
  key_name   = "deployer-one"
  public_key = tls_private_key.this.public_key_openssh
}

resource "aws_key_pair" "ssh_connecter" {
    key_name   = "deployer-key" 
    public_key = tls_private_key.this.public_key_openssh
}


resource "aws_instance" "api-ec2" {
    ami = "ami-072056ff9d3689e7b"
    instance_type = "t2.micro"
    security_groups = [aws_security_group.api_security_terra.name]
    key_name = aws_key_pair.ssh_connecter.key_name
    tags = {
        Name = "api instance"
    }
}
