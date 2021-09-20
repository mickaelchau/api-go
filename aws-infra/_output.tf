output "security_groups_id" {
    value = aws_security_group.api_security_terra.id
}

output "ec2_id" {
    value = aws_instance.api_ec2.id
}

output "ec2_keys_id" {
    value = aws_key_pair.ec2_keys.id
} 
