variable "service_serving_inbound_rules" {
  type = object({
    authorized_addresses = list(string)
    description          = string
    port                 = number
  })

  default = {
    authorized_addresses = ["0.0.0.0/0"]
    description          = "Serving service"
    port                 = 8000
  }
}

variable "ssh_inbound_rules" {
  type = object({
    authorized_addresses = list(string)
    description          = string
    port                 = number
  })

  default = {
    authorized_addresses = ["0.0.0.0/0"]
    description          = "SSH service"
    port                 = 22
  }
}

variable "https_inbound_rules" {
  type = object({
    authorized_addresses = list(string)
    description          = string
    port                 = number
  })

  default = {
    authorized_addresses = ["0.0.0.0/0"]
    description          = "HTTPS service"
    port                 = 443
  }
}

variable "outbound_rules" {
  type = object({
    authorized_addresses = list(string)
    description          = string
    port                 = number
  })
  
  default = {
    authorized_addresses = ["0.0.0.0/0"]
    description          = "Service comunicate with every ports"
    port                 = 0
  }
}

variable "service_name" {
  default = "Pokédex API REST"
}

variable "keys_type" {
  default = "RSA"
}

variable "instance_ami" {
  default = "ami-072056ff9d3689e7b"
}

variable "instance_type" {
  default = "t2.micro"
}
