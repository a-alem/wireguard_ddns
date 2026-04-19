// VPC
resource "aws_vpc" "vpn_onprem_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "kfupm-vpn-onprem-vpc"
    project = "master"
  }
}

resource "aws_vpc" "vpn_remote_vpc" {
  cidr_block = "10.10.0.0/16"

  tags = {
    Name = "kfupm-vpn-remote-vpc"
    project = "master"
  }
}

// Gateways
resource "aws_internet_gateway" "vpn_onprem_igw" {
  vpc_id = aws_vpc.vpn_onprem_vpc.id

  tags = {
    Name    = "kfupm-vpn-onprem-igw"
    project = "master"
  }
}

resource "aws_internet_gateway" "vpn_remote_igw" {
  vpc_id = aws_vpc.vpn_remote_vpc.id

  tags = {
    Name    = "kfupm-vpn-remote-igw"
    project = "master"
  }
}

// Routing tables
resource "aws_route_table" "vpn_onprem_public_rt" {
  vpc_id = aws_vpc.vpn_onprem_vpc.id

  tags = {
    Name    = "kfupm-vpn-onprem-public-rt"
    project = "master"
  }
}

resource "aws_route_table" "vpn_remote_public_rt" {
  vpc_id = aws_vpc.vpn_remote_vpc.id

  tags = {
    Name    = "kfupm-vpn-remote-public-rt"
    project = "master"
  }
}

// Default routes
resource "aws_route" "vpn_onprem_public_internet_route" {
  route_table_id         = aws_route_table.vpn_onprem_public_rt.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.vpn_onprem_igw.id
}

resource "aws_route" "vpn_remote_public_internet_route" {
  route_table_id         = aws_route_table.vpn_remote_public_rt.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.vpn_remote_igw.id
}

// Route table associations
resource "aws_route_table_association" "vpn_onprem_public_assoc" {
  subnet_id      = aws_subnet.vpn_onprem_subnet_public.id
  route_table_id = aws_route_table.vpn_onprem_public_rt.id
}

resource "aws_route_table_association" "vpn_remote_public_assoc" {
  subnet_id      = aws_subnet.vpn_remote_subnet_public.id
  route_table_id = aws_route_table.vpn_remote_public_rt.id
}

// Subnets
resource "aws_subnet" "vpn_onprem_subnet_public" {
  vpc_id     = aws_vpc.vpn_onprem_vpc.id
  availability_zone = var.av_zone
  map_public_ip_on_launch = true # This sets a subnet as public
  cidr_block = "10.0.0.0/20"

  tags = {
    Name = "kfupm-vpn-onprem-subnet-public1-eu-central-1a"
    project = "master"
  }
}

resource "aws_subnet" "vpn_onprem_subnet_private" {
  vpc_id     = aws_vpc.vpn_onprem_vpc.id
  availability_zone = var.av_zone
  cidr_block = "10.0.128.0/20"

  tags = {
    Name = "kfupm-vpn-onprem-subnet-private1-eu-central-1a"
    project = "master"
  }
}

resource "aws_subnet" "vpn_remote_subnet_public" {
  vpc_id     = aws_vpc.vpn_remote_vpc.id
  availability_zone = var.av_zone
  map_public_ip_on_launch = true # This sets a subnet as public
  cidr_block = "10.10.0.0/20"

  tags = {
    Name = "kfupm-vpn-remote-subnet-public1-eu-central-1a"
  }
}

resource "aws_subnet" "vpn_remote_subnet_private" {
  vpc_id     = aws_vpc.vpn_remote_vpc.id
  availability_zone = var.av_zone
  cidr_block = "10.10.128.0/20"

  tags = {
    Name = "kfupm-vpn-remote-subnet-private1-eu-central-1a"
  }
}

// Security groups
// For onprem subnet:
resource "aws_security_group" "allow_ssh_wireguard_onprem" {
  name        = "allow_ssh_wireguard_onprem"
  description = "Allow SSH and wireguard inbound traffic and all outbound traffic for onprem subnet"
  vpc_id      = aws_vpc.vpn_onprem_vpc.id

  tags = {
    Name = "allow_ssh_wireguard_onprem"
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_ssh_onprem" {
  security_group_id = aws_security_group.allow_ssh_wireguard_onprem.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  ip_protocol       = "tcp"
  to_port           = 22
}

resource "aws_vpc_security_group_ingress_rule" "allow_wireguard_onprem" {
  security_group_id = aws_security_group.allow_ssh_wireguard_onprem.id
  cidr_ipv4        = "0.0.0.0/0"
  from_port         = 51820
  ip_protocol       = "udp"
  to_port           = 51820
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4_onprem" {
  security_group_id = aws_security_group.allow_ssh_wireguard_onprem.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1" # all ports
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv6_onprem" {
  security_group_id = aws_security_group.allow_ssh_wireguard_onprem.id
  cidr_ipv6         = "::/0"
  ip_protocol       = "-1" # all ports
}

// For remote subnet:
resource "aws_security_group" "allow_ssh_wireguard_remote" {
  name        = "allow_ssh_wireguard_remote"
  description = "Allow SSH and wireguard inbound traffic and all outbound traffic for remote subnet"
  vpc_id      = aws_vpc.vpn_remote_vpc.id

  tags = {
    Name = "allow_ssh_wireguard_remote"
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_ssh_remote" {
  security_group_id = aws_security_group.allow_ssh_wireguard_remote.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  ip_protocol       = "tcp"
  to_port           = 22
}

resource "aws_vpc_security_group_ingress_rule" "allow_wireguard_remote" {
  security_group_id = aws_security_group.allow_ssh_wireguard_remote.id
  cidr_ipv4        = "0.0.0.0/0"
  from_port        = 51820
  ip_protocol      = "udp"
  to_port          = 51820
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4_remote" {
  security_group_id = aws_security_group.allow_ssh_wireguard_remote.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1" # all ports
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv6_remote" {
  security_group_id = aws_security_group.allow_ssh_wireguard_remote.id
  cidr_ipv6         = "::/0"
  ip_protocol       = "-1" # all ports
}