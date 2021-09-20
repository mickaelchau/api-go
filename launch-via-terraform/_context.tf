locals {
  application_name = "Pokédex API REST"
  keys_type        = "RSA"
  instance_ami     = "ami-072056ff9d3689e7b"
  instance_type    = "t2.micro"
  inbound_rules = [
    {
      authorized_addresses = ["0.0.0.0/0"]
      description          = "Serving application"
      protocol             = "tcp"
      port                 = 8000
    },
    {
      authorized_addresses = ["0.0.0.0/0"]
      description          = "SSH service"
      protocol             = "tcp"
      port                 = 22
    },
    {
      authorized_addresses = ["0.0.0.0/0"]
      description          = "HTTPS service"
      protocol             = "tcp"
      port                 = 443
    }
  ]

  outbound_rules = [
    {
      authorized_addresses = ["0.0.0.0/0"]
      description          = "Service comunicate with every ports"
      protocol             = "-1"
      port                 = 0
    }
  ]
}
