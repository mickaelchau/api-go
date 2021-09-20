resource "aws_security_group" "api_security_terra" {
  dynamic "ingress" {
    for_each = local.inbound_rules
    content {
      cidr_blocks = ingress.value["authorized_addresses"]
      description = ingress.value["description"]
      protocol    = ingress.value["protocol"]
      from_port   = ingress.value["port"]
      to_port     = ingress.value["port"]
    }    
  }
  dynamic "egress" {
    for_each = local.outbound_rules
    content {
      cidr_blocks = egress.value["authorized_addresses"]
      description = egress.value["description"]
      protocol    = egress.value["protocol"]
      from_port   = egress.value["port"]
      to_port     = egress.value["port"]
    }    
  }

  tags = {
    Name = "${local.application_name} security rules"
  }
}
