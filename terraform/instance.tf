// Instances EC2
resource "aws_instance" "vpn_onprem_server" {
  ami = "ami-0c1f1e2efc1f37d2e"
  instance_type = "t3.nano"
  availability_zone = "me-south-1a"
  subnet_id = aws_subnet.vpn_onprem_subnet_public.id

  tags = {
    Name = "vpn-terminal-server-onprem"
  }
}

resource "aws_instance" "vpn_remote_server" {
  ami = "ami-0c1f1e2efc1f37d2e"
  instance_type = "t3.nano"
  availability_zone = "me-south-1a"
  subnet_id = aws_subnet.vpn_remote_subnet_public.id

  tags = {
    Name = "vpn-terminal-server-remote"
  }
}