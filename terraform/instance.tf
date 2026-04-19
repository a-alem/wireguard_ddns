// Key pairs
resource "aws_key_pair" "wireguard-onprem" {
  key_name   = "kfupm-wireguard-onprem-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIF+1ZPJ0dJwtjeU8OG1vAfQn4XXnWrWar59maeloQB5o onprem wireguard"
}

resource "aws_key_pair" "wireguard-remote" {
  key_name   = "kfupm-wireguard-remote-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAINFR/W1JRUUIEXIOI8PVBFN5Tt80ozTJ1Atof76kZsc0 remote wireguard"
}

// Instances EC2
resource "aws_instance" "vpn_onprem_server" {
  ami = "ami-0281b0943230d40d1"
  instance_type = "t3.nano"
  availability_zone = var.av_zone
  subnet_id = aws_subnet.vpn_onprem_subnet_public.id
  key_name = aws_key_pair.wireguard-onprem.key_name
  vpc_security_group_ids = [
    aws_security_group.allow_ssh_wireguard_onprem.id
  ]

  tags = {
    Name = "vpn-terminal-server-onprem"
  }
}

resource "aws_instance" "vpn_remote_server" {
  ami = "ami-0281b0943230d40d1"
  instance_type = "t3.nano"
  availability_zone = var.av_zone
  subnet_id = aws_subnet.vpn_remote_subnet_public.id
  key_name = aws_key_pair.wireguard-remote.key_name
  vpc_security_group_ids = [
    aws_security_group.allow_ssh_wireguard_remote.id
  ]

  tags = {
    Name = "vpn-terminal-server-remote"
  }
}